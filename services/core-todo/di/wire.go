//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"google.golang.org/grpc"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/todo"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	todousecase "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos"
)

func InitializeApp(cfg *config.Config) (*grpc.Server, func(), error) {
	wire.Build(
		persistence.NewDatabase,

		persistence.WireSet,

		persistence.NewTodoCommandsGateway,
		wire.Bind(new(gateway.TodoCommandsGateway), new(*persistence.TodoCommandsGateway)),

		// Usecases
		todousecase.WireSet,

		// Handler + gRPC server
		todo.NewTodoHandler,
		todo.ProvideGRPCServer,
	)
	return nil, nil, nil
}
