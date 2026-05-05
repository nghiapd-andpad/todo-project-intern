package persistence

import (
	"github.com/google/wire"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
)

var WireSet = wire.NewSet(
	NewDatabase,

	NewTodoCommandsGateway,
	wire.Bind(new(gateway.TodoCommandsGateway), new(*TodoCommandsGateway)),

	NewTodoQueriesGateway,
	wire.Bind(new(gateway.TodoQueriesGateway), new(*TodoQueriesGateway)),

	NewTodoListCommandsGateway,
	wire.Bind(new(gateway.TodoListCommandsGateway), new(*TodoListCommandsGateway)),

	NewTodoListQueriesGateway,
	wire.Bind(new(gateway.TodoListQueriesGateway), new(*TodoListQueriesGateway)),
)
