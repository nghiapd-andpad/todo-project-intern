package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
)

func setupTopology(conn *amqp.Connection, cfg *config.Config) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbitmq open channel: %w", err)
	}
	defer ch.Close()

	// todo.events exchange
	if err := ch.ExchangeDeclare(
		cfg.RabbitMQTodoExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("declare todo exchange: %w", err)
	}

	// todo.notification.queue
	if _, err := ch.QueueDeclare(
		cfg.RabbitMQNotificationQueue,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("declare notification queue: %w", err)
	}

	if err := ch.QueueBind(
		cfg.RabbitMQNotificationQueue,
		cfg.RabbitMQNotificationRoutingKey, // "todo.assigned"
		cfg.RabbitMQTodoExchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("bind notification queue: %w", err)
	}

	return nil
}
