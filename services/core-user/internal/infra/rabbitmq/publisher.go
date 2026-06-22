package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
)

type Publisher struct {
	ch       *amqp.Channel
	exchange string
}

func NewPublisher(conn *amqp.Connection, cfg *config.Config) (*Publisher, func(), error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("rabbitmq open publisher channel: %w", err)
	}
	if err := ch.Confirm(false); err != nil {
		_ = ch.Close()
		return nil, nil, fmt.Errorf("rabbitmq set confirm mode: %w", err)
	}

	cleanup := func() { _ = ch.Close() }

	return &Publisher{
		ch:       ch,
		exchange: cfg.RabbitMQTodoExchange,
	}, cleanup, nil
}

var _ gateway.EventPublisher = (*Publisher)(nil)

func (p *Publisher) Publish(ctx context.Context, routingKey string, payload []byte) error {
	confirms := p.ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	if err := p.ch.PublishWithContext(
		ctx,
		p.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         payload,
		},
	); err != nil {
		return fmt.Errorf("rabbitmq publish: %w", err)
	}

	select {
	case confirm, ok := <-confirms:
		if !ok {
			return fmt.Errorf("rabbitmq confirm channel closed")
		}
		if !confirm.Ack {
			return fmt.Errorf("rabbitmq broker nacked message (routing_key=%s)", routingKey)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("rabbitmq publish confirm timeout: %w", ctx.Err())
	}
}
