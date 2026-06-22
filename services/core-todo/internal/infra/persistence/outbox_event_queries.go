package persistence

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/output"
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

func (g *OutboxEventQueriesGateway) FindClaimable(ctx context.Context, in *input.FindClaimableOutboxEvents) ([]*output.OutboxEvent, error) {
	conn := connFromContext(ctx, g.db)

	stuckBefore := time.Now().UTC().Add(-in.StuckThreshold)

	var models []*model.OutboxEvent

	if err := conn.
		Clauses(clause.Locking{
			Strength: "UPDATE",
			Options:  "SKIP LOCKED",
		}).
		Where(
			"(status IN (?, ?) OR (status = ? AND updated_at < ?)) AND retry_count < ?",
			string(output.OutboxEventStatusPending),
			string(output.OutboxEventStatusFailed),
			string(output.OutboxEventStatusProcessing),
			stuckBefore,
			in.MaxRetry,
		).
		Order("created_at ASC").
		Limit(in.BatchSize).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("db find claimable outbox events: %w", err)
	}

	return mapper.OutboxEventsToOutput(models), nil
}
