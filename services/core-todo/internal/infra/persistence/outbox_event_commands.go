package persistence

import (
	"context"
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	gatewayoutput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/output"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/model"
)

type OutboxEventCommandsGateway struct {
	db *gorm.DB
}

func NewOutboxEventCommandsGateway(db *gorm.DB) *OutboxEventCommandsGateway {
	return &OutboxEventCommandsGateway{db: db}
}

var _ gateway.OutboxEventCommandsGateway = (*OutboxEventCommandsGateway)(nil)

func (g *OutboxEventCommandsGateway) Create(ctx context.Context, in *gatewayinput.CreateOutboxEvent) error {
	conn := connFromContext(ctx, g.db)

	m := &model.OutboxEvent{
		EventName:  in.EventName,
		RoutingKey: in.RoutingKey,
		Payload:    datatypes.JSON(in.Payload),
		Status:     string(gatewayoutput.OutboxEventStatusPending),
	}

	if err := conn.Create(m).Error; err != nil {
		return fmt.Errorf("db create outbox event: %w", err)
	}

	return nil
}

func (g *OutboxEventCommandsGateway) MarkPublished(ctx context.Context, id int64) error {
	conn := connFromContext(ctx, g.db)

	now := time.Now().UTC()
	if err := conn.Model(&model.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       string(gatewayoutput.OutboxEventStatusSent),
			"published_at": now,
			"last_error":   nil,
		}).Error; err != nil {
		return fmt.Errorf("db mark outbox event published: %w", err)
	}

	return nil
}

func (g *OutboxEventCommandsGateway) MarkFailed(ctx context.Context, id int64, errMsg string) error {
	conn := connFromContext(ctx, g.db)

	if err := conn.Model(&model.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":      string(gatewayoutput.OutboxEventStatusFailed),
			"retry_count": gorm.Expr("retry_count + 1"),
			"last_error":  errMsg,
		}).Error; err != nil {
		return fmt.Errorf("db mark outbox event failed: %w", err)
	}

	return nil
}
