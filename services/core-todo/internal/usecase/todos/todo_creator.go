package todos

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoCreator interface {
	Create(ctx context.Context, in *input.TodoCreator) (*output.TodoCreator, error)
}

type todoCreator struct {
	todoCommandsGateway gateway.TodoCommandsGateway
}

func NewTodoCreator(todoCommandsGateway gateway.TodoCommandsGateway) TodoCreator {
	return &todoCreator{
		todoCommandsGateway: todoCommandsGateway,
	}
}

func (s *todoCreator) Create(ctx context.Context, in *input.TodoCreator) (*output.TodoCreator, error) {
	todo := &entity.Todo{
		Title:       in.Title,
		Description: in.Description,
		Status:      entity.TodoStatusPending,
	}

	todo, err := s.todoCommandsGateway.Create(ctx, todo)
	if err != nil {
		return nil, err
	}

	return &output.TodoCreator{
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
