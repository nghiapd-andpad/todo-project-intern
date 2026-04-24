package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
)

type TodoUpdater interface {
	UpdateTodoList(ctx context.Context, name string, displayName *string) (*entity.TodoList, error)
	UpdateTodo(ctx context.Context, name string, input gateway.UpdateTodoInput) (*entity.Todo, error)
}

type todoUpdater struct {
	todoGateway gateway.TodoGateway
}

func NewTodoUpdater(todoGateway gateway.TodoGateway) TodoUpdater {
	return &todoUpdater{todoGateway: todoGateway}
}

func (u *todoUpdater) UpdateTodoList(ctx context.Context, name string, displayName *string) (*entity.TodoList, error) {
	result, err := u.todoGateway.UpdateTodoList(ctx, name, displayName)
	if err != nil {
		return nil, fmt.Errorf("todoUpdater.UpdateTodoList: %w", err)
	}

	return result, nil
}

func (u *todoUpdater) UpdateTodo(ctx context.Context, name string, input gateway.UpdateTodoInput) (*entity.Todo, error) {
	result, err := u.todoGateway.UpdateTodo(ctx, name, input)
	if err != nil {
		return nil, fmt.Errorf("todoUpdater.UpdateTodo: %w", err)
	}

	return result, nil
}
