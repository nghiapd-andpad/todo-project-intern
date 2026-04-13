package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/user-service/internal/domain/entity"
)

type UserCommandsGateway interface {
	Create(ctx context.Context, user *entity.User) error
}
