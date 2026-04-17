//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/nghiapd-andpad/todo-project-intern/gateways/todo-bff/internal/handler/graph"
	"github.com/nghiapd-andpad/todo-project-intern/gateways/todo-bff/internal/infra/grpc_client"
	"github.com/nghiapd-andpad/todo-project-intern/gateways/todo-bff/internal/usecase/auth"
	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
)

func InitializeResolver(userSvc userv1.UserServiceClient) *graph.Resolver {
	wire.Build(
		grpc_client.ProviderSet,
		auth.ProviderSet,
		graph.NewResolver,
	)
	return nil
}
