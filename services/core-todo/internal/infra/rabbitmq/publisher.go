package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
)

type Publisher struct {
	conn     *amqp.Connection
	exchange string
}

func NewPublisher(conn *amqp.Connection, cfg *config.Config) *Publisher {
	return &Publisher{
		conn:     conn,
		exchange: cfg.RabbitMQExchange,
	}
}

var _ gateway.EventPublisher = (*Publisher)(nil)

func (p *Publisher) Publish(ctx context.Context, routingKey string, payload []byte) error {

	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.PublishWithContext(
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
	)
}
