package rabbitmq

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/event"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
)

type EmailConsumer struct {
	conn            *amqp.Connection
	emailSender     gateway.EmailSender
	processedEvents gateway.ProcessedEventGateway
	transactor      gateway.Transactor
	cfg             *config.Config
	logger          *zap.Logger
}

func NewEmailConsumer(
	conn *amqp.Connection,
	emailSender gateway.EmailSender,
	processedEvents gateway.ProcessedEventGateway,
	transactor gateway.Transactor,
	cfg *config.Config,
	zapLogger *zap.Logger,
) *EmailConsumer {
	if zapLogger == nil {
		zapLogger = zap.NewNop()
	}
	return &EmailConsumer{
		conn:            conn,
		emailSender:     emailSender,
		processedEvents: processedEvents,
		transactor:      transactor,
		cfg:             cfg,
		logger:          zapLogger.With(zap.String("component", "email_consumer")),
	}
}

func (c *EmailConsumer) Start(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("email consumer: open channel: %w", err)
	}

	if err := ch.Qos(10, 0, false); err != nil {
		_ = ch.Close()
		return fmt.Errorf("email consumer: set qos: %w", err)
	}

	deliveries, err := ch.Consume(
		"email.queue",
		"email-consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		return fmt.Errorf("email consumer: start consume: %w", err)
	}

	c.logger.Info("email consumer started", zap.String("queue", "email.queue"))

	go func() {
		defer func() {
			_ = ch.Close()
			c.logger.Info("email consumer stopped")
		}()

		var wg sync.WaitGroup
		workers := c.cfg.RabbitMQEmailWorkers
		if workers <= 0 {
			workers = 1
		}

		c.logger.Info("starting email consumer worker pool", zap.Int("workers", workers))

		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				c.logger.Debug("email consumer worker started", zap.Int("worker_id", workerID))
				for {
					select {
					case <-ctx.Done():
						return
					case d, ok := <-deliveries:
						if !ok {
							return
						}
						c.handle(ctx, d)
					}
				}
			}(i)
		}

		wg.Wait()
	}()

	return nil
}

func (c *EmailConsumer) handle(ctx context.Context, d amqp.Delivery) {
	logger := c.logger.With(
		zap.String("routing_key", d.RoutingKey),
		zap.Uint64("delivery_tag", d.DeliveryTag),
	)

	// Hash payload for technical duplicate check
	hashBytes := sha256.Sum256(d.Body)
	eventHash := hex.EncodeToString(hashBytes[:])

	var p event.NotificationEmailRequested
	if err := json.Unmarshal(d.Body, &p); err != nil {
		logger.Error("email consumer: unmarshal failed — discarding",
			zap.Error(err),
			zap.ByteString("body", d.Body),
		)
		_ = d.Ack(false)
		return
	}

	var isDuplicate bool
	err := c.transactor.Transaction(ctx, func(txCtx context.Context) error {
		ok, err := c.processedEvents.TryRecord(txCtx, eventHash, "email_consumer")
		if err != nil {
			return err
		}
		if !ok {
			isDuplicate = true
		}
		return nil
	})

	if err != nil {
		logger.Error("email consumer: database transaction failed — requeuing",
			zap.Error(err),
		)
		if getXDeathCount(d.Headers) >= 3 {
			_ = d.Nack(false, false)
			return
		}
		_ = d.Nack(false, true)
		return
	}

	if isDuplicate {
		_ = d.Ack(false)
		logger.Warn("email consumer: duplicate event detected via technical hash — skipped",
			zap.String("event_hash", eventHash),
		)
		return
	}

	logger.Info("email consumer: sending email",
		zap.String("to", p.Email),
		zap.String("subject", p.Subject),
	)

	if err := c.emailSender.Send(ctx, p.Email, p.Subject, p.Content); err != nil {
		logger.Error("email consumer: send email failed — requeuing",
			zap.String("to", p.Email),
			zap.Error(err),
		)
		if getXDeathCount(d.Headers) >= 3 {
			_ = d.Nack(false, false)
			return
		}
		_ = d.Nack(false, true)
		return
	}

	_ = d.Ack(false)
	logger.Info("email consumer: email sent successfully", zap.String("to", p.Email))
}
