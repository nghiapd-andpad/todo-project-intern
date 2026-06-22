package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/event"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
)

type EmailConsumer struct {
	conn        *amqp.Connection
	emailSender gateway.EmailSender
	cfg         *config.Config
	logger      *zap.Logger
}

func NewEmailConsumer(
	conn *amqp.Connection,
	emailSender gateway.EmailSender,
	cfg *config.Config,
	zapLogger *zap.Logger,
) *EmailConsumer {
	if zapLogger == nil {
		zapLogger = zap.NewNop()
	}
	return &EmailConsumer{
		conn:        conn,
		emailSender: emailSender,
		cfg:         cfg,
		logger:      zapLogger.With(zap.String("component", "email_consumer")),
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

		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-deliveries:
				if !ok {
					c.logger.Error("email consumer: delivery channel closed")
					return
				}
				c.handle(ctx, d)
			}
		}
	}()

	return nil
}

func (c *EmailConsumer) handle(ctx context.Context, d amqp.Delivery) {
	logger := c.logger.With(
		zap.String("routing_key", d.RoutingKey),
		zap.Uint64("delivery_tag", d.DeliveryTag),
	)

	var p event.NotificationEmailRequested
	if err := json.Unmarshal(d.Body, &p); err != nil {
		logger.Error("email consumer: unmarshal failed — discarding",
			zap.Error(err),
			zap.ByteString("body", d.Body),
		)
		_ = d.Ack(false)
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
