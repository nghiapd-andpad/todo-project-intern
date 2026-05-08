package auth

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	inputgateway "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/auth/input"
)

type Registerer struct {
	authGateway gateway.AuthGateway
}

func NewRegisterer(authGateway gateway.AuthGateway) *Registerer {
	return &Registerer{authGateway: authGateway}
}

func (u *Registerer) Register(ctx context.Context, input input.RegisterInput) (*entity.User, error) {
	user, err := u.authGateway.Register(ctx, inputgateway.RegisterInput{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("Registerer.Register: %w", err)
	}

	return user, nil
}
