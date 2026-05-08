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
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/infra/grpcclient/mapper"
)

type UserGateway struct {
	client userv1.UserServiceClient
}

var _ gateway.UserGateway = (*UserGateway)(nil)

func NewUserGateway(cfg *config.Config) (*UserGateway, func(), error) {
	conn, err := grpc.Dial(
		cfg.UserServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(auth.UnaryClientInterceptor()),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("dial user service: %w", err)
	}

	client := userv1.NewUserServiceClient(conn)

	return &UserGateway{client: client}, func() { conn.Close() }, nil
}

func (g *UserGateway) GetByID(ctx context.Context, id string) (*entity.User, error) {
	resp, err := g.client.GetUser(ctx, &userv1.GetUserRequest{Id: id})
	if err != nil {
		return nil, fmt.Errorf("UserGateway.GetByID: %w", err)
	}

	return mapper.UserFromPb(resp.User), nil
}

func (g *UserGateway) GetByIDs(ctx context.Context, ids []string) ([]*entity.User, error) {
	resp, err := g.client.BatchGetUsers(ctx, &userv1.BatchGetUsersRequest{
		Ids: ids,
	})
	if err != nil {
		return nil, fmt.Errorf("UserGateway.GetByIDs: %w", err)
	}

	users := make([]*entity.User, len(resp.GetUsers()))
	for i, pbUser := range resp.GetUsers() {
		users[i] = mapper.UserFromPb(pbUser)
	}

	return users, nil
}

func (g *UserGateway) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	resp, err := g.client.GetUserByUsername(ctx, &userv1.GetUserByUsernameRequest{
		Username: username,
	})
	if err != nil {
		return nil, fmt.Errorf("UserGateway.GetByUsername: %w", err)
	}

	return mapper.UserFromPb(resp.User), nil
}

func (g *UserGateway) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	resp, err := g.client.GetUserByEmail(ctx, &userv1.GetUserByEmailRequest{
		Email: email,
	})
	if err != nil {
		return nil, fmt.Errorf("UserGateway.GetByEmail: %w", err)
	}

	return mapper.UserFromPb(resp.User), nil
}
