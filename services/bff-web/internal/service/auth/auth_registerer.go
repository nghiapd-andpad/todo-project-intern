package auth

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	inputgateway "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/output"
)

type Registerer struct {
	authGateway gateway.AuthGateway
}

func NewRegisterer(authGateway gateway.AuthGateway) *Registerer {
	return &Registerer{authGateway: authGateway}
}

func (u *Registerer) Register(ctx context.Context, input *input.RegisterInput) (*output.RegisterOutput, error) {
	user, err := u.authGateway.Register(ctx, inputgateway.RegisterInput{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("Registerer.Register: %w", err)
	}

	return &output.RegisterOutput{
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
