package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
)

type TodoDeleter struct {
	todoGateway gateway.TodoGateway
}

func NewTodoDeleter(todoGateway gateway.TodoGateway) *TodoDeleter {
	return &TodoDeleter{todoGateway: todoGateway}
}

func (u *TodoDeleter) DeleteTodoList(ctx context.Context, name string) error {
	if err := u.todoGateway.DeleteTodoList(ctx, name); err != nil {
		return fmt.Errorf("TodoDeleter.DeleteTodoList: %w", err)
	}

	return nil
}

func (u *TodoDeleter) DeleteTodo(ctx context.Context, name string) error {
	if err := u.todoGateway.DeleteTodo(ctx, name); err != nil {
		return fmt.Errorf("TodoDeleter.DeleteTodo: %w", err)
	}

	return nil
}
