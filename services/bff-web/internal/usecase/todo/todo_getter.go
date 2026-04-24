package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
)

type TodoGetter interface {
	GetTodoList(ctx context.Context, name string) (*entity.TodoList, error)
	GetTodo(ctx context.Context, name string) (*entity.Todo, error)
}

type todoGetter struct {
	todoGateway gateway.TodoGateway
}

func NewTodoGetter(todoGateway gateway.TodoGateway) TodoGetter {
	return &todoGetter{todoGateway: todoGateway}
}

func (u *todoGetter) GetTodoList(ctx context.Context, name string) (*entity.TodoList, error) {
	result, err := u.todoGateway.GetTodoList(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("todoGetter.GetTodoList: %w", err)
	}

	return result, nil
}

func (u *todoGetter) GetTodo(ctx context.Context, name string) (*entity.Todo, error) {
	result, err := u.todoGateway.GetTodo(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("todoGetter.GetTodo: %w", err)
	}

	return result, nil
}
