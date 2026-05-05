// Package todo provides use cases related to todo management, such as creating todo lists and todos.
package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/todo/input"
)

type TodoCreator interface {
	CreateTodoList(ctx context.Context, input input.CreateTodoListInput) (*entity.TodoList, error)
	CreateTodo(ctx context.Context, parent string, input input.CreateTodoInput) (*entity.Todo, error)
}

type todoCreator struct {
	todoGateway gateway.TodoGateway
}

func NewTodoCreator(todoGateway gateway.TodoGateway) TodoCreator {
	return &todoCreator{todoGateway: todoGateway}
}

func (u *todoCreator) CreateTodoList(ctx context.Context, input input.CreateTodoListInput) (*entity.TodoList, error) {
	if input.DisplayName == "" {
		return nil, entity.NewInvalidParameter("display_name is required")
	}

	result, err := u.todoGateway.CreateTodoList(ctx, gateway.CreateTodoListInput{
		Parent:      input.Parent,
		DisplayName: input.DisplayName,
	})
	if err != nil {
		return nil, fmt.Errorf("todoCreator.CreateTodoList: %w", err)
	}

	return result, nil
}

func (u *todoCreator) CreateTodo(ctx context.Context, parent string, input input.CreateTodoInput) (*entity.Todo, error) {
	if input.Title == "" {
		return nil, entity.NewInvalidParameter("title is required")
	}

	result, err := u.todoGateway.CreateTodo(ctx, parent, gateway.CreateTodoInput{
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
		AssigneeID:  input.AssigneeID,
	})
	if err != nil {
		return nil, fmt.Errorf("todoCreator.CreateTodo: %w", err)
	}

	return result, nil
}
