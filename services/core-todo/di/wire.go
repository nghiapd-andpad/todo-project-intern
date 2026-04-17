//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/todo"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos"
	"google.golang.org/grpc"
)

func InitializeApp() (*grpc.Server, func(), error) {
	wire.Build(
		config.New,
		persistence.NewDatabase,
		persistence.WireSet,
		todos.WireSet,
		todo.NewTodoHandler,
		todo.NewGRPCServer,
	)
	return nil, nil, nil
}
