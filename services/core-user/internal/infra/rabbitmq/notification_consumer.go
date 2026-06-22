package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway/input"
)

type NotificationConsumer struct {
	conn      *amqp.Connection
	notifCmds gateway.NotificationCommandsGateway
	cfg       *config.Config
	logger    *zap.Logger
}

func NewNotificationConsumer(
	conn *amqp.Connection,
	notifCmds gateway.NotificationCommandsGateway,
	cfg *config.Config,
	zapLogger *zap.Logger,
) *NotificationConsumer {
	if zapLogger == nil {
		zapLogger = zap.NewNop()
	}
	return &NotificationConsumer{
		conn:      conn,
		notifCmds: notifCmds,
		cfg:       cfg,
		logger:    zapLogger.With(zap.String("component", "notification_consumer")),
	}
}

func (c *NotificationConsumer) Start(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("notification consumer: open channel: %w", err)
	}

	if err := ch.Qos(c.cfg.RabbitMQNotificationPrefetchCount, 0, false); err != nil {
		_ = ch.Close()
		return fmt.Errorf("notification consumer: set qos: %w", err)
	}

	deliveries, err := ch.Consume(
		c.cfg.RabbitMQNotificationQueue,
		"notification-consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		return fmt.Errorf("notification consumer: start consume: %w", err)
	}

	c.logger.Info("notification consumer started",
		zap.String("queue", c.cfg.RabbitMQNotificationQueue))

	go func() {
		defer func() {
			_ = ch.Close()
			c.logger.Info("notification consumer stopped")
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-deliveries:
				if !ok {
					c.logger.Error("notification consumer: delivery channel closed")
					return
				}
				c.handle(ctx, d)
			}
		}
	}()

	return nil
}

// todoAssignedPayload mirrors event.TodoAssigned from core-todo.
type todoAssignedPayload struct {
	TodoID     int64     `json:"todo_id"`
	TodoListID int64     `json:"todo_list_id"`
	ActorID    int64     `json:"actor_id"`
	AssigneeID int64     `json:"assignee_id"`
	Title      string    `json:"title"`
	OccurredOn time.Time `json:"occurred_at"`
}

func (c *NotificationConsumer) handle(ctx context.Context, d amqp.Delivery) {
	logger := c.logger.With(
		zap.String("routing_key", d.RoutingKey),
		zap.Uint64("delivery_tag", d.DeliveryTag),
	)

	var p todoAssignedPayload
	if err := json.Unmarshal(d.Body, &p); err != nil {
		logger.Error("notification consumer: unmarshal failed — discarding",
			zap.Error(err),
			zap.ByteString("body", d.Body),
		)
		_ = d.Ack(false)
		return
	}

	if err := c.notifCmds.Create(ctx, &gatewayinput.CreateNotification{
		ReceiverID: p.AssigneeID,
		Title:      "You have been assigned a new todo.",
		Content:    fmt.Sprintf("Todo \"%s\" assigned a new todo.", p.Title),
	}); err != nil {
		logger.Error("notification consumer: insert notification failed — requeuing",
			zap.Int64("assignee_id", p.AssigneeID),
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

	logger.Info("notification consumer: notification created",
		zap.Int64("receiver_id", p.AssigneeID),
		zap.Int64("todo_id", p.TodoID),
	)
}

// getXDeathCount reads the x-death header RabbitMQ attaches on each NACK+requeue.
func getXDeathCount(headers amqp.Table) int {
	deaths, ok := headers["x-death"].([]interface{})
	if !ok || len(deaths) == 0 {
		return 0
	}
	entry, ok := deaths[0].(amqp.Table)
	if !ok {
		return 0
	}
	count, _ := entry["count"].(int64)
	return int(count)
}
