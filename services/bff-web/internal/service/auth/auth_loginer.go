// Package auth provides use cases related to user authentication, such as login and registration.
package auth

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	inputgateway "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/output"
)

type Loginer struct {
	authGateway gateway.AuthGateway
}

func NewLoginer(authGateway gateway.AuthGateway) *Loginer {
	return &Loginer{authGateway: authGateway}
}

func (u *Loginer) Login(ctx context.Context, input *input.LoginInput) (*output.LoginOutput, error) {
	data, err := u.authGateway.Login(ctx, inputgateway.LoginInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("Loginer.Login: %w", err)
	}

	return &output.LoginOutput{
		AccessToken: data.AccessToken,
		User:        data.User,
	}, nil
}
