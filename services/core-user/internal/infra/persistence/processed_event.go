package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence/model"
)

type ProcessedEventGateway struct {
	db *gorm.DB
}

func NewProcessedEventGateway(db *gorm.DB) *ProcessedEventGateway {
	return &ProcessedEventGateway{db: db}
}

var _ gateway.ProcessedEventGateway = (*ProcessedEventGateway)(nil)

func (g *ProcessedEventGateway) TryRecord(ctx context.Context, hash string, consumerKey string) (bool, error) {
	conn := connFromContext(ctx, g.db)

	m := &model.ProcessedEvent{
		EventHash:   hash,
		ConsumerKey: consumerKey,
	}

	result := conn.Clauses(clause.OnConflict{DoNothing: true}).Create(m)
	if err := result.Error; err != nil {
		return false, fmt.Errorf("db record processed event: %w", err)
	}

	return result.RowsAffected > 0, nil
}
