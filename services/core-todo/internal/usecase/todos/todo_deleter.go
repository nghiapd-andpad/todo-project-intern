package todos

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoDeleter struct {
	todoCommandsGateway gateway.TodoCommandsGateway
	todoQueriesGateway  gateway.TodoQueriesGateway
}

func NewTodoDeleter(todoCommandsGateway gateway.TodoCommandsGateway, todoQueriesGateway gateway.TodoQueriesGateway) *TodoDeleter {
	return &TodoDeleter{todoCommandsGateway: todoCommandsGateway, todoQueriesGateway: todoQueriesGateway}
}

func (s *TodoDeleter) Delete(ctx context.Context, in *input.TodoDeleter) (*output.TodoDeleter, error) {
	todo, err := s.todoQueriesGateway.Get(ctx, in.ID)
	if err != nil {
		return nil, fmt.Errorf("TodoDeleter.Get: %w", err)
	}
	if todo == nil {
		return nil, entity.NewNotFound("todo not found").
			WithDetail("todo_id", fmt.Sprintf("%d", in.ID))
	}

	// delete todo
	if err := s.todoCommandsGateway.Delete(ctx, in.ID); err != nil {
		return nil, fmt.Errorf("TodoDeleter.Delete: %w", err)
	}

	return &output.TodoDeleter{}, nil
}
