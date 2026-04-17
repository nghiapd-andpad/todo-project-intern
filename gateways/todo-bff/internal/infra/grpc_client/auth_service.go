package grpc_client

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/gateways/todo-bff/internal/domain"
	"github.com/nghiapd-andpad/todo-project-intern/gateways/todo-bff/internal/gateway"
	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
)

type authServiceClient struct {
	client userv1.UserServiceClient
}

func NewAuthServiceClient(client userv1.UserServiceClient) gateway.AuthGateway {
	return &authServiceClient{client: client}
}

func (s *authServiceClient) Register(ctx context.Context, username, password, email string) (*domain.User, error) {
	resp, err := s.client.Register(ctx, &userv1.RegisterRequest{
		Username: username,
		Password: password,
		Email:    email,
	})
	if err != nil {
		return nil, fmt.Errorf("grpc register: %w", err)
	}

	return &domain.User{
		ID:       fmt.Sprintf("users/%v", resp.User.Id),
		Username: resp.User.Username,
		Email:    resp.User.Email,
	}, nil
}

func (s *authServiceClient) Login(ctx context.Context, username, password string) (string, *domain.User, error) {
	resp, err := s.client.Login(ctx, &userv1.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", nil, fmt.Errorf("grpc login: %w", err)
	}

	user := &domain.User{
		ID:       fmt.Sprintf("users/%v", resp.User.Id),
		Username: resp.User.Username,
	}

	return resp.AccessToken, user, nil
}
