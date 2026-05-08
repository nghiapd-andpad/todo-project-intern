package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
)

type TodoGetter struct {
	todoGateway gateway.TodoGateway
}

func NewTodoGetter(todoGateway gateway.TodoGateway) *TodoGetter {
	return &TodoGetter{todoGateway: todoGateway}
}

func (u *TodoGetter) GetTodoList(ctx context.Context, name string) (*entity.TodoList, error) {
	result, err := u.todoGateway.GetTodoList(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("TodoGetter.GetTodoList: %w", err)
	}

	return result, nil
}

func (u *TodoGetter) GetTodo(ctx context.Context, name string) (*entity.Todo, error) {
	result, err := u.todoGateway.GetTodo(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("TodoGetter.GetTodo: %w", err)
	}

	return result, nil
}
