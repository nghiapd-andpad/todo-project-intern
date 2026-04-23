package auth

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
)

type Loginer interface {
	Login(ctx context.Context, username, password string) (string, *entity.User, error)
}

type loginer struct {
	authGateway gateway.AuthGateway
}

func NewLoginer(authGateway gateway.AuthGateway) Loginer {
	return &loginer{authGateway: authGateway}
}

func (u *loginer) Login(ctx context.Context, username, password string) (string, *entity.User, error) {
	if username == "" || password == "" {
		return "", nil, entity.NewInvalidParameter("username and password are required")
	}
	token, user, err := u.authGateway.Login(ctx, username, password)
	if err != nil {
		return "", nil, fmt.Errorf("loginer.Login: %w", err)
	}
	return token, user, nil
}
