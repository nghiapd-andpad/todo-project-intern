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
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/event"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway/input"
)

type mockTransactor struct{}

func (m *mockTransactor) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return fn(ctx)
}

type mockUserQueries struct {
	getUserFn func(ctx context.Context, id entity.UserID) (*entity.User, error)
}

func (m *mockUserQueries) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return nil, nil
}
func (m *mockUserQueries) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}
func (m *mockUserQueries) GetByID(ctx context.Context, id entity.UserID) (*entity.User, error) {
	return m.getUserFn(ctx, id)
}
func (m *mockUserQueries) GetByIDs(ctx context.Context, ids []entity.UserID) ([]*entity.User, error) {
	return nil, nil
}

type mockNotifCommands struct {
	createFn func(ctx context.Context, in *gatewayinput.CreateNotification) (int64, int64, error)
}

func (m *mockNotifCommands) Create(ctx context.Context, in *gatewayinput.CreateNotification) (int64, int64, error) {
	return m.createFn(ctx, in)
}

type mockOutboxCommands struct {
	createFn func(ctx context.Context, in *gatewayinput.CreateOutboxEvent) error
}

func (m *mockOutboxCommands) Create(ctx context.Context, in *gatewayinput.CreateOutboxEvent) error {
	return m.createFn(ctx, in)
}
func (m *mockOutboxCommands) MarkProcessing(ctx context.Context, ids []int64) error {
	return nil
}
func (m *mockOutboxCommands) MarkPublished(ctx context.Context, id int64) error {
	return nil
}
func (m *mockOutboxCommands) MarkFailed(ctx context.Context, id int64, errMsg string) error {
	return nil
}
func (m *mockOutboxCommands) MarkDead(ctx context.Context, id int64, errMsg string) error {
	return nil
}

type mockProcessedEvents struct {
	tryRecordFn func(ctx context.Context, hash string, consumerKey string) (bool, error)
}

func (m *mockProcessedEvents) TryRecord(ctx context.Context, hash string, consumerKey string) (bool, error) {
	if m.tryRecordFn != nil {
		return m.tryRecordFn(ctx, hash, consumerKey)
	}
	return true, nil
}

func TestNotificationConsumer_Handle(t *testing.T) {
	cfg := &config.Config{
		RabbitMQNotificationQueue: "notification.queue",
	}

	t.Run("successfully processes first-time notification and creates outbox event", func(t *testing.T) {
		userQueries := &mockUserQueries{
			getUserFn: func(ctx context.Context, id entity.UserID) (*entity.User, error) {
				assert.Equal(t, entity.UserID(42), id)
				return &entity.User{
					ID:       42,
					Username: "john_doe",
					Email:    "john@example.com",
				}, nil
			},
		}

		notifCreated := false
		notifCmds := &mockNotifCommands{
			createFn: func(ctx context.Context, in *gatewayinput.CreateNotification) (int64, int64, error) {
				assert.Equal(t, int64(42), in.ReceiverID)
				assert.Equal(t, "todo", in.ResourceType)
				assert.Equal(t, int64(101), in.ResourceID)
				notifCreated = true
				return 7, 1, nil // ID = 7, RowsAffected = 1 (inserted)
			},
		}

		outboxCreated := false
		outboxCmds := &mockOutboxCommands{
			createFn: func(ctx context.Context, in *gatewayinput.CreateOutboxEvent) error {
				assert.Equal(t, "notification.email.requested", in.EventName)
				assert.Equal(t, "notification.email.requested", in.RoutingKey)

				var emailEvent event.NotificationEmailRequested
				err := json.Unmarshal(in.Payload, &emailEvent)
				assert.NoError(t, err)
				assert.Equal(t, int64(7), emailEvent.NotificationID)
				assert.Equal(t, "john@example.com", emailEvent.Email)
				assert.Contains(t, emailEvent.Content, "john_doe")
				outboxCreated = true
				return nil
			},
		}

		processedEvents := &mockProcessedEvents{
			tryRecordFn: func(ctx context.Context, hash string, consumerKey string) (bool, error) {
				assert.NotEmpty(t, hash)
				assert.Equal(t, "notification_consumer", consumerKey)
				return true, nil
			},
		}

		consumer := &NotificationConsumer{
			conn:            nil,
			notifCmds:       notifCmds,
			outboxCmds:      outboxCmds,
			userQueries:     userQueries,
			processedEvents: processedEvents,
			transactor:      &mockTransactor{},
			cfg:             cfg,
			logger:          zap.NewNop(),
		}

		payload := todoAssignedPayload{
			TodoID:     101,
			TodoListID: 202,
			ActorID:    1,
			AssigneeID: 42,
			Title:      "Complete intern project",
			OccurredOn: time.Now(),
		}
		body, err := json.Marshal(payload)
		assert.NoError(t, err)

		d := amqp.Delivery{
			RoutingKey: "todo.assigned",
			Body:       body,
		}

		consumer.handle(context.Background(), d)

		assert.True(t, notifCreated)
		assert.True(t, outboxCreated)
	})

	t.Run("ignores technical duplicate event", func(t *testing.T) {
		userQueries := &mockUserQueries{
			getUserFn: func(ctx context.Context, id entity.UserID) (*entity.User, error) {
				return &entity.User{
					ID:       42,
					Username: "john_doe",
					Email:    "john@example.com",
				}, nil
			},
		}

		notifCmds := &mockNotifCommands{
			createFn: func(ctx context.Context, in *gatewayinput.CreateNotification) (int64, int64, error) {
				t.Error("notification should not be created for duplicate events")
				return 0, 0, nil
			},
		}

		outboxCmds := &mockOutboxCommands{
			createFn: func(ctx context.Context, in *gatewayinput.CreateOutboxEvent) error {
				t.Error("outbox event should not be created for duplicate events")
				return nil
			},
		}

		processedEvents := &mockProcessedEvents{
			tryRecordFn: func(ctx context.Context, hash string, consumerKey string) (bool, error) {
				return false, nil // Already processed (duplicate)
			},
		}

		consumer := &NotificationConsumer{
			conn:            nil,
			notifCmds:       notifCmds,
			outboxCmds:      outboxCmds,
			userQueries:     userQueries,
			processedEvents: processedEvents,
			transactor:      &mockTransactor{},
			cfg:             cfg,
			logger:          zap.NewNop(),
		}

		payload := todoAssignedPayload{
			TodoID:     101,
			TodoListID: 202,
			ActorID:    1,
			AssigneeID: 42,
			Title:      "Complete intern project",
			OccurredOn: time.Now(),
		}
		body, err := json.Marshal(payload)
		assert.NoError(t, err)

		d := amqp.Delivery{
			RoutingKey: "todo.assigned",
			Body:       body,
		}

		consumer.handle(context.Background(), d)
	})

	t.Run("handles assignee not found gracefully", func(t *testing.T) {
		userQueries := &mockUserQueries{
			getUserFn: func(ctx context.Context, id entity.UserID) (*entity.User, error) {
				return nil, nil // Not found
			},
		}

		consumer := &NotificationConsumer{
			conn:            nil,
			notifCmds:       nil,
			outboxCmds:      nil,
			userQueries:     userQueries,
			processedEvents: &mockProcessedEvents{},
			transactor:      &mockTransactor{},
			cfg:             cfg,
			logger:          zap.NewNop(),
		}

		payload := todoAssignedPayload{
			TodoID:     101,
			TodoListID: 202,
			ActorID:    1,
			AssigneeID: 42,
			Title:      "Complete intern project",
			OccurredOn: time.Now(),
		}
		body, err := json.Marshal(payload)
		assert.NoError(t, err)

		d := amqp.Delivery{
			RoutingKey: "todo.assigned",
			Body:       body,
		}

		consumer.handle(context.Background(), d)
	})

	t.Run("handles user queries failure with retry", func(t *testing.T) {
		userQueries := &mockUserQueries{
			getUserFn: func(ctx context.Context, id entity.UserID) (*entity.User, error) {
				return nil, errors.New("database connection issue")
			},
		}

		consumer := &NotificationConsumer{
			conn:            nil,
			notifCmds:       nil,
			outboxCmds:      nil,
			userQueries:     userQueries,
			processedEvents: &mockProcessedEvents{},
			transactor:      &mockTransactor{},
			cfg:             cfg,
			logger:          zap.NewNop(),
		}

		payload := todoAssignedPayload{
			TodoID:     101,
			TodoListID: 202,
			ActorID:    1,
			AssigneeID: 42,
			Title:      "Complete intern project",
			OccurredOn: time.Now(),
		}
		body, err := json.Marshal(payload)
		assert.NoError(t, err)

		d := amqp.Delivery{
			RoutingKey: "todo.assigned",
			Body:       body,
		}

		consumer.handle(context.Background(), d)
	})
}
