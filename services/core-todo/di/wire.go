//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"google.golang.org/grpc"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	todoHandler "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/todo"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos"
)

func InitializeApp(cfg *config.Config) (*grpc.Server, func(), error) {
	wire.Build(
		// Infra
		persistence.NewDatabase,

		persistence.NewTodoCommandsGateway,
		wire.Bind(new(gateway.TodoCommandsGateway), new(*persistence.TodoCommandsGateway)),

		persistence.NewTodoQueriesGateway,
		wire.Bind(new(gateway.TodoQueriesGateway), new(*persistence.TodoQueriesGateway)),

		persistence.NewTodoListCommandsGateway,
		wire.Bind(new(gateway.TodoListCommandsGateway), new(*persistence.TodoListCommandsGateway)),

		persistence.NewTodoListQueriesGateway,
		wire.Bind(new(gateway.TodoListQueriesGateway), new(*persistence.TodoListQueriesGateway)),

		// Todo Usecase
		todos.NewTodoCreator,
		todos.NewTodoGetter,
		todos.NewTodoLister,
		todos.NewTodoUpdater,
		todos.NewTodoDeleter,
		todos.NewTodoListCreator,
		todos.NewTodoListGetter,
		todos.NewTodoListLister,
		todos.NewTodoListUpdater,
		todos.NewTodoListDeleter,

		// Handler
		todoHandler.NewTodoHandler,
		todoHandler.NewGRPCServer,
	)
	return nil, nil, nil
}
