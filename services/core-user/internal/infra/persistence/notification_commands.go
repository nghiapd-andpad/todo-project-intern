package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/persistence/model"
)

type NotificationCommandsGateway struct {
	db *gorm.DB
}

func NewNotificationCommandsGateway(db *gorm.DB) *NotificationCommandsGateway {
	return &NotificationCommandsGateway{db: db}
}

var _ gateway.NotificationCommandsGateway = (*NotificationCommandsGateway)(nil)

func (g *NotificationCommandsGateway) Create(ctx context.Context, in *gatewayinput.CreateNotification) (int64, int64, error) {
	conn := connFromContext(ctx, g.db)

	m := &model.Notification{
		ReceiverID:   in.ReceiverID,
		ResourceType: in.ResourceType,
		ResourceID:   in.ResourceID,
		EventName:    in.EventName,
		OccurredAt:   in.OccurredAt,
		Title:        in.Title,
		Content:      in.Content,
	}

	result := conn.Clauses(clause.OnConflict{DoNothing: true}).Create(m)
	if err := result.Error; err != nil {
		return 0, 0, fmt.Errorf("db create notification: %w", err)
	}

	return m.ID, result.RowsAffected, nil
}
