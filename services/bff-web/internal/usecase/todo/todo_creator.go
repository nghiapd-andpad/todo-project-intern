// Package todo provides use cases related to todo management, such as creating todo lists and todos.
package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
)

type TodoCreator interface {
	CreateTodoList(ctx context.Context, parent string, displayName string) (*entity.TodoList, error)
	CreateTodo(ctx context.Context, parent string, input gateway.CreateTodoInput) (*entity.Todo, error)
}

type todoCreator struct {
	todoGateway gateway.TodoGateway
}

func NewTodoCreator(todoGateway gateway.TodoGateway) TodoCreator {
	return &todoCreator{todoGateway: todoGateway}
}

func (u *todoCreator) CreateTodoList(ctx context.Context, parent string, displayName string) (*entity.TodoList, error) {
	if displayName == "" {
		return nil, entity.NewInvalidParameter("display_name is required")
	}

	result, err := u.todoGateway.CreateTodoList(ctx, parent, displayName)
	if err != nil {
		return nil, fmt.Errorf("todoCreator.CreateTodoList: %w", err)
	}

	return result, nil
}

func (u *todoCreator) CreateTodo(ctx context.Context, parent string, input gateway.CreateTodoInput) (*entity.Todo, error) {
	if input.Title == "" {
		return nil, entity.NewInvalidParameter("title is required")
	}

	result, err := u.todoGateway.CreateTodo(ctx, parent, input)
	if err != nil {
		return nil, fmt.Errorf("todoCreator.CreateTodo: %w", err)
	}

	return result, nil
}
