package persistence

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	gatewayoutput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/output"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/model"
)

type IdempotencyGateway struct {
	db *gorm.DB
}

func NewIdempotencyGateway(db *gorm.DB) *IdempotencyGateway {
	return &IdempotencyGateway{db: db}
}

var _ gateway.IdempotencyGateway = (*IdempotencyGateway)(nil)

func (g *IdempotencyGateway) Find(
	ctx context.Context,
	userID entity.UserID,
	operation string,
	key string,
) (*gatewayoutput.IdempotencyRecord, error) {
	conn := connFromContext(ctx, g.db)

	var m model.IdempotencyKey
	err := conn.
		Where("user_id = ? AND operation = ? AND idempotency_key = ?", int64(userID), operation, key).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("db find idempotency key: %w", err)
	}

	return mapper.IdempotencyRecordFromModel(&m), nil
}

func (g *IdempotencyGateway) CreateProcessing(
	ctx context.Context,
	in *gatewayinput.CreateIdempotencyRecord,
) (*gatewayoutput.IdempotencyRecord, error) {
	conn := connFromContext(ctx, g.db)

	m := &model.IdempotencyKey{
		UserID:         int64(in.UserID),
		Operation:      in.Operation,
		IdempotencyKey: in.IdempotencyKey,
		RequestHash:    in.RequestHash,
		Status:         string(gatewayoutput.IdempotencyStatusProcessing),
		ExpiresAt:      in.ExpiresAt,
	}

	if err := conn.Create(m).Error; err != nil {
		return nil, fmt.Errorf("db create idempotency key: %w", err)
	}

	return mapper.IdempotencyRecordFromModel(m), nil
}

func (g *IdempotencyGateway) MarkCompleted(ctx context.Context, id int64, resourceID int64) error {
	conn := connFromContext(ctx, g.db)

	if err := conn.Model(&model.IdempotencyKey{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":      string(gatewayoutput.IdempotencyStatusCompleted),
			"resource_id": resourceID,
		}).Error; err != nil {
		return fmt.Errorf("db mark idempotency completed: %w", err)
	}

	return nil
}
