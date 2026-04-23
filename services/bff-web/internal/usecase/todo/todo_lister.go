package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
)

type TodoLister interface {
	ListTodoLists(ctx context.Context, parent string, opts gateway.ListTodoListsOptions) (*gateway.TodoListPage, error)
	ListTodos(ctx context.Context, parent string, opts gateway.ListTodosOptions) (*gateway.TodoPage, error)
}

type todoLister struct {
	todoGateway gateway.TodoGateway
}

func NewTodoLister(todoGateway gateway.TodoGateway) TodoLister {
	return &todoLister{todoGateway: todoGateway}
}

func (u *todoLister) ListTodoLists(ctx context.Context, parent string, opts gateway.ListTodoListsOptions) (*gateway.TodoListPage, error) {
	result, err := u.todoGateway.ListTodoLists(ctx, parent, opts)
	if err != nil {
		return nil, fmt.Errorf("todoLister.ListTodoLists: %w", err)
	}
	return result, nil
}

func (u *todoLister) ListTodos(ctx context.Context, parent string, opts gateway.ListTodosOptions) (*gateway.TodoPage, error) {
	result, err := u.todoGateway.ListTodos(ctx, parent, opts)
	if err != nil {
		return nil, fmt.Errorf("todoLister.ListTodos: %w", err)
	}
	return result, nil
}
