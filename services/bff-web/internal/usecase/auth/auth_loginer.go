// Package auth provides use cases related to user authentication, such as login and registration.
package auth

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/auth/input"
)

type Loginer interface {
	Login(ctx context.Context, input input.LoginInput) (string, *entity.User, error)
}

type loginer struct {
	authGateway gateway.AuthGateway
}

func NewLoginer(authGateway gateway.AuthGateway) Loginer {
	return &loginer{authGateway: authGateway}
}

func (u *loginer) Login(ctx context.Context, input input.LoginInput) (string, *entity.User, error) {
	if input.Username == "" || input.Password == "" {
		return "", nil, entity.NewInvalidParameter("username and password are required")
	}

	token, user, err := u.authGateway.Login(ctx, gateway.LoginInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		return "", nil, fmt.Errorf("loginer.Login: %w", err)
	}

	return token, user, nil
}
