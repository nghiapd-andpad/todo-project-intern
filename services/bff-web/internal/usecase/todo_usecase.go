package usecase

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/output"
)

type TodoCreator interface {
	CreateTodoList(ctx context.Context, in *input.CreateTodoListInput) (*output.TodoListOutput, error)
	CreateTodo(ctx context.Context, in *input.CreateTodoInput) (*output.TodoOutput, error)
}

type TodoGetter interface {
	GetTodoList(ctx context.Context, name string) (*output.TodoListOutput, error)
	GetTodo(ctx context.Context, name string) (*output.TodoOutput, error)
}

type TodoLister interface {
	ListTodoLists(ctx context.Context, parent string, opts *input.ListTodoListsOptions) (*output.TodoListPage, error)
	ListTodos(ctx context.Context, parent string, opts *input.ListTodosOptions) (*output.TodoPage, error)
}

type TodoUpdater interface {
	UpdateTodoList(ctx context.Context, in *input.UpdateTodoListInput) (*output.TodoListOutput, error)
	UpdateTodo(ctx context.Context, in *input.UpdateTodoInput) (*output.TodoOutput, error)
}

type TodoDeleter interface {
	DeleteTodoList(ctx context.Context, name string) error
	DeleteTodo(ctx context.Context, name string) error
}
