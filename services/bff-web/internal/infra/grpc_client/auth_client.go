package grpc_client

import (
	"context"
	"fmt"

	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/gateway"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, err
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
		return "", nil, err
	}

	user := &domain.User{
		ID:       fmt.Sprintf("users/%v", resp.User.Id),
		Username: resp.User.Username,
	}

	return resp.AccessToken, user, nil
}

func (s *authServiceClient) VerifyToken(ctx context.Context, token string) (string, []string, error) {
	resp, err := s.client.VerifyToken(ctx, &userv1.VerifyTokenRequest{
		AccessToken: token,
	})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.Unauthenticated:
			return "", nil, domain.ErrUnauthorized
		case codes.PermissionDenied:
			return "", nil, domain.ErrForbidden
		default:
			return "", nil, domain.ErrInternal
		}
	}

	return resp.UserId, resp.Roles, nil
}
