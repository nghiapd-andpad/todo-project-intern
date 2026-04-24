package todos

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoDeleter interface {
	Delete(ctx context.Context, in *input.TodoDeleter) (*output.TodoDeleter, error)
}

type todoDeleter struct {
	todoCommandsGateway gateway.TodoCommandsGateway
	todoQueriesGateway  gateway.TodoQueriesGateway
}

func NewTodoDeleter(todoCommandsGateway gateway.TodoCommandsGateway, todoQueriesGateway gateway.TodoQueriesGateway) TodoDeleter {
	return &todoDeleter{todoCommandsGateway: todoCommandsGateway, todoQueriesGateway: todoQueriesGateway}
}

func (s *todoDeleter) Delete(ctx context.Context, in *input.TodoDeleter) (*output.TodoDeleter, error) {
	todo, err := s.todoQueriesGateway.Get(ctx, in.ID)
	if err != nil {
		return nil, fmt.Errorf("todoDeleter.Get: %w", err)
	}
	if todo == nil {
		return nil, entity.NewNotFound("todo not found").
			WithDetail("todo_id", fmt.Sprintf("%d", in.ID))
	}

	// delete todo
	if err := s.todoCommandsGateway.Delete(ctx, in.ID); err != nil {
		return nil, fmt.Errorf("todoCommandsGateway.Delete: %w", err)
	}

	return &output.TodoDeleter{}, nil
}
