package usecase

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/output"
)

type AuthRegisterer interface {
	Register(ctx context.Context, in *input.RegisterInput) (*output.RegisterOutput, error)
}

type AuthLoginer interface {
	Login(ctx context.Context, in *input.LoginInput) (*output.LoginOutput, error)
}
