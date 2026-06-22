package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/event"
)

type mockEmailSender struct {
	sendFn func(ctx context.Context, to, subject, content string) error
}

func (m *mockEmailSender) Send(ctx context.Context, to, subject, content string) error {
	if m.sendFn != nil {
		return m.sendFn(ctx, to, subject, content)
	}
	return nil
}

func TestEmailConsumer_Handle(t *testing.T) {
	cfg := &config.Config{}

	t.Run("successfully sends email and acks delivery", func(t *testing.T) {
		emailSent := false
		emailSender := &mockEmailSender{
			sendFn: func(ctx context.Context, to, subject, content string) error {
				assert.Equal(t, "john@example.com", to)
				assert.Equal(t, "Hello Subject", subject)
				assert.Equal(t, "Hello Content", content)
				emailSent = true
				return nil
			},
		}

		processedEvents := &mockProcessedEvents{
			tryRecordFn: func(ctx context.Context, hash string, consumerKey string) (bool, error) {
				assert.NotEmpty(t, hash)
				assert.Equal(t, "email_consumer", consumerKey)
				return true, nil
			},
		}

		consumer := &EmailConsumer{
			conn:            nil,
			emailSender:     emailSender,
			processedEvents: processedEvents,
			transactor:      &mockTransactor{},
			cfg:             cfg,
			logger:          zap.NewNop(),
		}

		emailEvent := event.NotificationEmailRequested{
			NotificationID: 7,
			ReceiverID:     42,
			Email:          "john@example.com",
			Subject:        "Hello Subject",
			Content:        "Hello Content",
			OccurredOn:     time.Now(),
		}
		body, err := json.Marshal(emailEvent)
		assert.NoError(t, err)

		d := amqp.Delivery{
			RoutingKey: "notification.email.requested",
			Body:       body,
		}

		consumer.handle(context.Background(), d)

		assert.True(t, emailSent)
	})

	t.Run("skips email sending if event is duplicate", func(t *testing.T) {
		emailSender := &mockEmailSender{
			sendFn: func(ctx context.Context, to, subject, content string) error {
				t.Error("email should not be sent for duplicate events")
				return nil
			},
		}

		processedEvents := &mockProcessedEvents{
			tryRecordFn: func(ctx context.Context, hash string, consumerKey string) (bool, error) {
				return false, nil // already processed
			},
		}

		consumer := &EmailConsumer{
			conn:            nil,
			emailSender:     emailSender,
			processedEvents: processedEvents,
			transactor:      &mockTransactor{},
			cfg:             cfg,
			logger:          zap.NewNop(),
		}

		emailEvent := event.NotificationEmailRequested{
			NotificationID: 7,
			ReceiverID:     42,
			Email:          "john@example.com",
			Subject:        "Hello Subject",
			Content:        "Hello Content",
			OccurredOn:     time.Now(),
		}
		body, err := json.Marshal(emailEvent)
		assert.NoError(t, err)

		d := amqp.Delivery{
			RoutingKey: "notification.email.requested",
			Body:       body,
		}

		consumer.handle(context.Background(), d)
	})

	t.Run("handles email sending error gracefully", func(t *testing.T) {
		emailSender := &mockEmailSender{
			sendFn: func(ctx context.Context, to, subject, content string) error {
				return errors.New("smtp connection failed")
			},
		}

		consumer := &EmailConsumer{
			conn:            nil,
			emailSender:     emailSender,
			processedEvents: &mockProcessedEvents{},
			transactor:      &mockTransactor{},
			cfg:             cfg,
			logger:          zap.NewNop(),
		}

		emailEvent := event.NotificationEmailRequested{
			NotificationID: 7,
			ReceiverID:     42,
			Email:          "john@example.com",
			Subject:        "Hello Subject",
			Content:        "Hello Content",
			OccurredOn:     time.Now(),
		}
		body, err := json.Marshal(emailEvent)
		assert.NoError(t, err)

		d := amqp.Delivery{
			RoutingKey: "notification.email.requested",
			Body:       body,
		}

		consumer.handle(context.Background(), d)
	})
}
