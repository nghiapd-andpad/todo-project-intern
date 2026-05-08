package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/output"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/infra/grpcclient/mapper"
)

type AuthGateway struct {
	client userv1.UserServiceClient
}

var _ gateway.AuthGateway = (*AuthGateway)(nil)

func NewAuthGateway(cfg *config.Config) (*AuthGateway, func(), error) {
	conn, err := grpc.Dial(
		cfg.UserServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(auth.UnaryClientInterceptor()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("dial user service: %w", err)
	}
	client := userv1.NewUserServiceClient(conn)
	return &AuthGateway{client: client}, func() { conn.Close() }, nil
}

func (g *AuthGateway) Register(ctx context.Context, input input.RegisterInput) (*entity.User, error) {
	resp, err := g.client.Register(ctx, &userv1.RegisterRequest{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
	})
	if err != nil {
		return nil, err
	}

	return mapper.UserFromPb(resp.User), nil
}

func (g *AuthGateway) Login(ctx context.Context, input input.LoginInput) (*output.LoginOutput, error) {
	resp, err := g.client.Login(ctx, &userv1.LoginRequest{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	return &output.LoginOutput{
		AccessToken: resp.AccessToken,
		User:        mapper.UserFromPb(resp.User),
	}, nil
}
