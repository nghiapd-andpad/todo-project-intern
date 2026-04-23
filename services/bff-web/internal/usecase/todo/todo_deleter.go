package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
)

type TodoDeleter interface {
	DeleteTodoList(ctx context.Context, name string) error
	DeleteTodo(ctx context.Context, name string) error
}

type todoDeleter struct {
	todoGateway gateway.TodoGateway
}

func NewTodoDeleter(todoGateway gateway.TodoGateway) TodoDeleter {
	return &todoDeleter{todoGateway: todoGateway}
}

func (u *todoDeleter) DeleteTodoList(ctx context.Context, name string) error {
	if err := u.todoGateway.DeleteTodoList(ctx, name); err != nil {
		return fmt.Errorf("todoDeleter.DeleteTodoList: %w", err)
	}
	return nil
}

func (u *todoDeleter) DeleteTodo(ctx context.Context, name string) error {
	if err := u.todoGateway.DeleteTodo(ctx, name); err != nil {
		return fmt.Errorf("todoDeleter.DeleteTodo: %w", err)
	}
	return nil
}
