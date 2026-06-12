package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
)

func NewConnection(cfg *config.Config) (*amqp.Connection, func(), error) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.RabbitMQUser,
		cfg.RabbitMQPassword,
		cfg.RabbitMQHost,
		cfg.RabbitMQPort,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("rabbitmq dial: %w", err)
	}

	if err := setupTopology(conn, cfg); err != nil {
		_ = conn.Close()
		return nil, nil, err
	}

	cleanup := func() {
		_ = conn.Close()
	}

	return conn, cleanup, nil
}
