//go:generate mockgen -destination=mock/notification_commands_mock.go -source=notification_commands.go -package mock

package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway/input"
)

type NotificationCommandsGateway interface {
	Create(ctx context.Context, in *input.CreateNotification) (id int64, rowsAffected int64, err error)
}
