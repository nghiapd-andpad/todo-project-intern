//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"google.golang.org/grpc"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/service"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase"
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
		service.NewTodoCreator,
		service.NewTodoGetter,
		service.NewTodoLister,
		service.NewTodoUpdater,
		service.NewTodoDeleter,

		service.NewTodoListCreator,
		service.NewTodoListGetter,
		service.NewTodoListLister,
		service.NewTodoListUpdater,
		service.NewTodoListDeleter,

		wire.Bind(new(usecase.TodoCreator), new(*service.TodoCreator)),
		wire.Bind(new(usecase.TodoGetter), new(*service.TodoGetter)),
		wire.Bind(new(usecase.TodoLister), new(*service.TodoLister)),
		wire.Bind(new(usecase.TodoUpdater), new(*service.TodoUpdater)),
		wire.Bind(new(usecase.TodoDeleter), new(*service.TodoDeleter)),

		wire.Bind(new(usecase.TodoListCreator), new(*service.TodoListCreator)),
		wire.Bind(new(usecase.TodoListGetter), new(*service.TodoListGetter)),
		wire.Bind(new(usecase.TodoListLister), new(*service.TodoListLister)),
		wire.Bind(new(usecase.TodoListUpdater), new(*service.TodoListUpdater)),
		wire.Bind(new(usecase.TodoListDeleter), new(*service.TodoListDeleter)),

		// HANDLER
		handler.NewTodoHandler,

		// SERVER
		handler.NewGRPCServer,
	)
	return nil, nil, nil
}
