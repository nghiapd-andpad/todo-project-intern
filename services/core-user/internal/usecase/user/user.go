package user

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user/output"
)

type UserCreator interface {
	Register(ctx context.Context, in *input.UserRegister) (*output.UserRegister, error)
}

type UserAuthenticator interface {
	Login(ctx context.Context, in *input.UserLogin) (*output.UserLogin, error)
	Verify(ctx context.Context, token string) (*output.VerifyTokenOutput, error)
}
