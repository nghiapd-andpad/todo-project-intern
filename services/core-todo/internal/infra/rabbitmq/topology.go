package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
)

func setupTopology(
	conn *amqp.Connection,
	cfg *config.Config,
) error {

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbitmq open channel: %w", err)
	}
	defer ch.Close()

	// Topic Exchange

	if err := ch.ExchangeDeclare(
		cfg.RabbitMQExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}

	// Audit Queue

	_, err = ch.QueueDeclare(
		cfg.RabbitMQAuditQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare audit queue: %w", err)
	}

	// Binding

	if err := ch.QueueBind(
		cfg.RabbitMQAuditQueue,
		cfg.RabbitMQAuditRoutingKey,
		cfg.RabbitMQExchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("bind audit queue: %w", err)
	}

	return nil
}
