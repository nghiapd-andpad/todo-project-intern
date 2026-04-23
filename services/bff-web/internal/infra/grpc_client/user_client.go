package grpc_client

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/infra/grpc_client/mapper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type userGateway struct {
	client userv1.UserServiceClient
}

func NewUserGateway(cfg *config.Config) (gateway.UserGateway, func(), error) {
	conn, err := grpc.Dial(
		cfg.UserServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(auth.UnaryClientInterceptor()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("dial user service: %w", err)
	}
	client := userv1.NewUserServiceClient(conn)
	return &userGateway{client: client}, func() { conn.Close() }, nil
}

func (g *userGateway) GetByID(ctx context.Context, id string) (*entity.User, error) {
	resp, err := g.client.GetUser(ctx, &userv1.GetUserRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("userGateway.GetByID: %w", err)
	}
	return mapper.UserFromPb(resp.User), nil
}

func (g *userGateway) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	resp, err := g.client.GetUserByUsername(ctx, &userv1.GetUserByUsernameRequest{
		Username: username,
	})
	if err != nil {
		return nil, fmt.Errorf("userGateway.GetByUsername: %w", err)
	}
	return mapper.UserFromPb(resp.User), nil
}

func (g *userGateway) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	resp, err := g.client.GetUserByEmail(ctx, &userv1.GetUserByEmailRequest{
		Email: email,
	})
	if err != nil {
		return nil, fmt.Errorf("userGateway.GetByEmail: %w", err)
	}
	return mapper.UserFromPb(resp.User), nil
}
