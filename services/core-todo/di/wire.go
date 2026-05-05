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
		// INFRASTRUCTURE
		persistence.NewDatabase,

		persistence.NewTodoCommandsGateway,
		persistence.NewTodoQueriesGateway,
		persistence.NewTodoListCommandsGateway,
		persistence.NewTodoListQueriesGateway,

		wire.Bind(new(gateway.TodoCommandsGateway), new(*persistence.TodoCommandsGateway)),
		wire.Bind(new(gateway.TodoQueriesGateway), new(*persistence.TodoQueriesGateway)),
		wire.Bind(new(gateway.TodoListCommandsGateway), new(*persistence.TodoListCommandsGateway)),
		wire.Bind(new(gateway.TodoListQueriesGateway), new(*persistence.TodoListQueriesGateway)),

		// USE CASE
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

		wire.Bind(new(todoHandler.TodoCreatorUsecase), new(*todos.TodoCreator)),
		wire.Bind(new(todoHandler.TodoGetterUsecase), new(*todos.TodoGetter)),
		wire.Bind(new(todoHandler.TodoListerUsecase), new(*todos.TodoLister)),
		wire.Bind(new(todoHandler.TodoUpdaterUsecase), new(*todos.TodoUpdater)),
		wire.Bind(new(todoHandler.TodoDeleterUsecase), new(*todos.TodoDeleter)),

		wire.Bind(new(todoHandler.TodoListCreatorUsecase), new(*todos.TodoListCreator)),
		wire.Bind(new(todoHandler.TodoListGetterUsecase), new(*todos.TodoListGetter)),
		wire.Bind(new(todoHandler.TodoListListerUsecase), new(*todos.TodoListLister)),
		wire.Bind(new(todoHandler.TodoListUpdaterUsecase), new(*todos.TodoListUpdater)),
		wire.Bind(new(todoHandler.TodoListDeleterUsecase), new(*todos.TodoListDeleter)),

		// HANDLER
		todoHandler.NewTodoHandler,

		// SERVER
		todoHandler.NewGRPCServer,
	)
	return nil, nil, nil
}
