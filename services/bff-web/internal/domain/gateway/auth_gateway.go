package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
)

type AuthGateway interface {
	Register(ctx context.Context, username, password, email string) (*entity.User, error)
	Login(ctx context.Context, username, password string) (string, *entity.User, error)
}
