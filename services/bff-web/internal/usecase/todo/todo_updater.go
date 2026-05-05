package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/todo/input"
)

type TodoUpdater interface {
	UpdateTodoList(ctx context.Context, input input.UpdateTodoListInput) (*entity.TodoList, error)
	UpdateTodo(ctx context.Context, name string, input input.UpdateTodoInput) (*entity.Todo, error)
}

type todoUpdater struct {
	todoGateway gateway.TodoGateway
}

func NewTodoUpdater(todoGateway gateway.TodoGateway) TodoUpdater {
	return &todoUpdater{todoGateway: todoGateway}
}

func (u *todoUpdater) UpdateTodoList(ctx context.Context, input input.UpdateTodoListInput) (*entity.TodoList, error) {
	result, err := u.todoGateway.UpdateTodoList(ctx, gateway.UpdateTodoListInput{
		Name:        input.Name,
		DisplayName: input.DisplayName,
	})
	if err != nil {
		return nil, fmt.Errorf("todoUpdater.UpdateTodoList: %w", err)
	}

	return result, nil
}

func (u *todoUpdater) UpdateTodo(ctx context.Context, name string, input input.UpdateTodoInput) (*entity.Todo, error) {
	result, err := u.todoGateway.UpdateTodo(ctx, name, gateway.UpdateTodoInput{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
		AssigneeID:  input.AssigneeID,
	})
	if err != nil {
		return nil, fmt.Errorf("todoUpdater.UpdateTodo: %w", err)
	}

	return result, nil
}
