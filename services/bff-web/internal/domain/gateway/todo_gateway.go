// Package gateway defines the interfaces for interacting with external services or data sources.
package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/output"
)

type TodoGateway interface {
	// TodoList operations
	GetTodoList(ctx context.Context, name string) (*entity.TodoList, error)
	ListTodoLists(ctx context.Context, parent string, opts input.ListTodoListsOptions) (*output.TodoListPage, error)
	CreateTodoList(ctx context.Context, input input.CreateTodoListInput) (*entity.TodoList, error)
	UpdateTodoList(ctx context.Context, input input.UpdateTodoListInput) (*entity.TodoList, error)
	DeleteTodoList(ctx context.Context, name string) error

	// Todo operations
	GetTodo(ctx context.Context, name string) (*entity.Todo, error)
	ListTodos(ctx context.Context, parent string, opts input.ListTodosOptions) (*output.TodoPage, error)
	CreateTodo(ctx context.Context, input input.CreateTodoInput) (*entity.Todo, error)
	UpdateTodo(ctx context.Context, input input.UpdateTodoInput) (*entity.Todo, error)
	DeleteTodo(ctx context.Context, name string) error
}
