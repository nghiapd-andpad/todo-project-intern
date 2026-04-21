//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/todo"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	todousecase "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos"
	"google.golang.org/grpc"
)

func InitializeApp(cfg *config.Config) (*grpc.Server, func(), error) {
	wire.Build(
		// Infrastructure
		persistence.NewDatabase,
		persistence.WireSet,

		// Usecases
		todousecase.WireSet,

		// Handler + gRPC server
		todo.NewTodoHandler,
		todo.NewGRPCServer,
	)
	return nil, nil, nil
}
