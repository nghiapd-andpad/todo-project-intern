package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway/input"
)

type NotificationCommandsGateway struct {
	db *gorm.DB
}

func NewNotificationCommandsGateway(db *gorm.DB) *NotificationCommandsGateway {
	return &NotificationCommandsGateway{db: db}
}

var _ gateway.NotificationCommandsGateway = (*NotificationCommandsGateway)(nil)

func (g *NotificationCommandsGateway) Create(ctx context.Context, in *gatewayinput.CreateNotification) error {
	if err := g.db.WithContext(ctx).Exec(
		`INSERT IGNORE INTO notifications
            (receiver_id, resource_type, resource_id, event_name, occurred_at,
             title, content, is_read, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, false, NOW(6), NOW(6))`,
		in.ReceiverID,
		in.ResourceType,
		in.ResourceID,
		in.EventName,
		in.OccurredAt,
		in.Title,
		in.Content,
	).Error; err != nil {
		return fmt.Errorf("db create notification: %w", err)
	}

	return nil
}
