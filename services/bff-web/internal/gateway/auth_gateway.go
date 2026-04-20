package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain"
)

type AuthGateway interface {
	Register(ctx context.Context, username, password, email string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (string, *domain.User, error)
	VerifyToken(ctx context.Context, token string) (string, []string, error)
}
