package todos

import (
	"context"
	"errors"

	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/usecase/todos/output"
)

type TodoUpdater interface {
	Update(ctx context.Context, in *input.TodoUpdater) (*output.TodoUpdater, error)
}

type todoUpdater struct {
	todoCommandsGateway gateway.TodoCommandsGateway
	todoQueriesGateway  gateway.TodoQueriesGateway
}

func NewTodoUpdater(todoCommandsGateway gateway.TodoCommandsGateway, todoQueriesGateway gateway.TodoQueriesGateway) TodoUpdater {
	return &todoUpdater{
		todoCommandsGateway: todoCommandsGateway,
		todoQueriesGateway:  todoQueriesGateway,
	}
}

func (s *todoUpdater) Update(ctx context.Context, in *input.TodoUpdater) (*output.TodoUpdater, error) {
	todo, err := s.todoQueriesGateway.Get(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if todo == nil {
		return nil, errors.New("todo not found")
	}

	if in.Title != nil {
		todo.Title = *in.Title
	}

	if in.Description != nil {
		todo.Description = in.Description
	}

	if in.Status != nil {
		todo.Status = *in.Status
	}

	if in.Priority != nil {
		todo.Priority = *in.Priority
	}

	todo, err = s.todoCommandsGateway.Update(ctx, todo)
	if err != nil {
		return nil, err
	}

	return &output.TodoUpdater{
		Todo: &entity.Todo{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      todo.Status,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		},
	}, nil
}
