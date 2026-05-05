package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
)

type RegisterInput struct {
	Username string
	Password string
	Email    string
}

type LoginInput struct {
	Username string
	Password string
}

type AuthGateway interface {
	Register(ctx context.Context, input RegisterInput) (*entity.User, error)
	Login(ctx context.Context, input LoginInput) (string, *entity.User, error)
}
