//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/handler/graph"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/infra/grpc_client"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/auth"
)

func InitializeResolver() (*graph.Resolver, func(), error) {
	wire.Build(
		grpc_client.ProviderSet,
		auth.ProviderSet,
		graph.NewResolver,
	)
	return nil, nil, nil
}
