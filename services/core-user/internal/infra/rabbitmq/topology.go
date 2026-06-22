package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
)

func setupTopology(conn *amqp.Connection, cfg *config.Config) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbitmq open channel: %w", err)
	}
	defer ch.Close()

	// Declare exchange

	// consume exchange
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

	// publish exchange
	if err := ch.ExchangeDeclare(
		cfg.RabbitMQUserExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// Declare notification queue
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

	// Bind: todo.assigned -> todo.notification.queue
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
