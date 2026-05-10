// Package todo provides use cases related to todo management, such as creating todo lists and todos.
package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	inputgateway "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/todo/input"
)

type TodoCreator struct {
	todoGateway gateway.TodoGateway
}

func NewTodoCreator(todoGateway gateway.TodoGateway) *TodoCreator {
	return &TodoCreator{todoGateway: todoGateway}
}

func (u *TodoCreator) CreateTodoList(ctx context.Context, input *input.CreateTodoListInput) (*entity.TodoList, error) {
	result, err := u.todoGateway.CreateTodoList(ctx, inputgateway.CreateTodoListInput{
		Parent:      input.Parent,
		DisplayName: input.DisplayName,
	})
	if err != nil {
		return nil, fmt.Errorf("TodoCreator.CreateTodoList: %w", err)
	}

	return result, nil
}

func (u *TodoCreator) CreateTodo(ctx context.Context, input *input.CreateTodoInput) (*entity.Todo, error) {
	result, err := u.todoGateway.CreateTodo(ctx, inputgateway.CreateTodoInput{
		Parent:      input.Parent,
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
		AssigneeID:  input.AssigneeID,
	})
	if err != nil {
		return nil, fmt.Errorf("TodoCreator.CreateTodo: %w", err)
	}

	return result, nil
}
