package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/output"
)

type AuthGateway interface {
	Register(ctx context.Context, input input.RegisterInput) (*entity.User, error)
	Login(ctx context.Context, input input.LoginInput) (*output.LoginOutput, error)
}
