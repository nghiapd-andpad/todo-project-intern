package auth

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
)

type Registerer interface {
	Register(ctx context.Context, username, password, email string) (*entity.User, error)
}

type registerer struct {
	authGateway gateway.AuthGateway
}

func NewRegisterer(authGateway gateway.AuthGateway) Registerer {
	return &registerer{authGateway: authGateway}
}

func (u *registerer) Register(ctx context.Context, username, password, email string) (*entity.User, error) {
	if username == "" || password == "" || email == "" {
		return nil, entity.NewInvalidParameter("username, password and email are required")
	}
	user, err := u.authGateway.Register(ctx, username, password, email)
	if err != nil {
		return nil, fmt.Errorf("registerer.Register: %w", err)
	}
	return user, nil
}
