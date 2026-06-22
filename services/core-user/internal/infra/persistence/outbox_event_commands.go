package persistence

import (
	"context"
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway/input"
	gatewayoutput "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway/output"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence/model"
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

func (g *OutboxEventCommandsGateway) MarkProcessing(ctx context.Context, ids []int64) error {
	conn := connFromContext(ctx, g.db)

	if len(ids) == 0 {
		return nil
	}

	if err := conn.Model(&model.OutboxEvent{}).
		Where("id IN ?", ids).
		Updates(map[string]any{
			"status":      string(gatewayoutput.OutboxEventStatusProcessing),
			"retry_count": gorm.Expr("retry_count + 1"),
		}).Error; err != nil {
		return fmt.Errorf("db mark outbox events processing: %w", err)
	}

	return nil
}

func (g *OutboxEventCommandsGateway) MarkPublished(ctx context.Context, id int64) error {
	now := time.Now().UTC()

	result := g.db.WithContext(ctx).
		Model(&model.OutboxEvent{}).
		Where("id = ? AND status = ?", id, string(gatewayoutput.OutboxEventStatusProcessing)).
		Updates(map[string]any{
			"status":       string(gatewayoutput.OutboxEventStatusPublished),
			"published_at": now,
			"last_error":   nil,
		})

	if result.Error != nil {
		return fmt.Errorf("db mark outbox event published: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return entity.NewAlreadyHandled(fmt.Sprintf("event %d already handled elsewhere", id))
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

func (g *OutboxEventCommandsGateway) MarkDead(ctx context.Context, id int64, errMsg string) error {
	conn := connFromContext(ctx, g.db)

	if err := conn.Model(&model.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     string(gatewayoutput.OutboxEventStatusDead),
			"last_error": errMsg,
		}).Error; err != nil {
		return fmt.Errorf("db mark outbox event dead: %w", err)
	}

	return nil
}
