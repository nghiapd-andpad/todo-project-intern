package auth

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/gateways/todo-bff/internal/domain"
	"github.com/nghiapd-andpad/todo-project-intern/gateways/todo-bff/internal/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/gateways/todo-bff/internal/usecase/auth/input"
	"github.com/nghiapd-andpad/todo-project-intern/gateways/todo-bff/internal/usecase/auth/output"
)

type AuthUseCase interface {
	Register(ctx context.Context, in *input.RegisterInput) (*domain.User, error)
	Login(ctx context.Context, in *input.LoginInput) (*output.LoginOutput, error)
}

type authUseCase struct {
	authGateway gateway.AuthGateway
}

func NewAuthUseCase(ag gateway.AuthGateway) AuthUseCase {
	return &authUseCase{authGateway: ag}
}

func (u *authUseCase) Register(ctx context.Context, in *input.RegisterInput) (*domain.User, error) {
	return u.authGateway.Register(ctx, in.Username, in.Password, in.Email)
}

func (u *authUseCase) Login(ctx context.Context, in *input.LoginInput) (*output.LoginOutput, error) {
	token, user, err := u.authGateway.Login(ctx, in.Username, in.Password)
	if err != nil {
		return nil, err
	}
	return &output.LoginOutput{
		AccessToken: token,
		User:        user,
	}, nil
}
