package auth

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/auth/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/auth/output"
)

type AuthUseCase interface {
	Register(ctx context.Context, in *input.RegisterInput) (*domain.User, error)
	Login(ctx context.Context, in *input.LoginInput) (*output.LoginOutput, error)
	Authenticate(ctx context.Context, token string) (string, []string, error)
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

func (u *authUseCase) Authenticate(ctx context.Context, token string) (string, []string, error) {
	if token == "" {
		return "", nil, domain.ErrUnauthorized
	}

	userID, roles, err := u.authGateway.VerifyToken(ctx, token)
	if err != nil {
		return "", nil, err
	}

	return userID, roles, nil
}
