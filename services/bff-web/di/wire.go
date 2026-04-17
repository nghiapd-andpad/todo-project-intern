//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/handler/graph"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/infra/grpc_client"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/auth"
)

func InitializeResolver(userSvc userv1.UserServiceClient) *graph.Resolver {
	wire.Build(
		grpc_client.ProviderSet,
		auth.ProviderSet,
		graph.NewResolver,
	)
	return nil
}
