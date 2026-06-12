package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	gatewayoutput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/output"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/model"
)

type OutboxEventQueriesGateway struct {
	db *gorm.DB
}

func NewOutboxEventQueriesGateway(db *gorm.DB) *OutboxEventQueriesGateway {
	return &OutboxEventQueriesGateway{db: db}
}

var _ gateway.OutboxEventQueriesGateway = (*OutboxEventQueriesGateway)(nil)

func (g *OutboxEventQueriesGateway) FindPending(
	ctx context.Context,
	in *gatewayinput.ListPendingOutboxEvents,
) ([]*gatewayoutput.OutboxEvent, error) {
	conn := connFromContext(ctx, g.db)

	limit := in.Limit
	if limit <= 0 {
		limit = 100
	}

	var models []*model.OutboxEvent
	if err := conn.
		Where("status = ?", string(gatewayoutput.OutboxEventStatusPending)).
		Order("created_at ASC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("db find pending outbox events: %w", err)
	}

	return mapper.OutboxEventsToOutput(models), nil
}
