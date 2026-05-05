package auth

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/auth/input"
)

type Registerer interface {
	Register(ctx context.Context, input input.RegisterInput) (*entity.User, error)
}

type registerer struct {
	authGateway gateway.AuthGateway
}

func NewRegisterer(authGateway gateway.AuthGateway) Registerer {
	return &registerer{authGateway: authGateway}
}

func (u *registerer) Register(ctx context.Context, input input.RegisterInput) (*entity.User, error) {
	if input.Username == "" || input.Password == "" || input.Email == "" {
		return nil, entity.NewInvalidParameter("username, password and email are required")
	}

	user, err := u.authGateway.Register(ctx, gateway.RegisterInput{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("registerer.Register: %w", err)
	}

	return user, nil
}
