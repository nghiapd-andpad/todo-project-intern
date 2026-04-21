package todos

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoDeleter interface {
	Delete(ctx context.Context, in *input.TodoDeleter) (*output.TodoDeleter, error)
}

type todoDeleter struct {
	todoCommandsGateway gateway.TodoCommandsGateway
}

func NewTodoDeleter(todoCommandsGateway gateway.TodoCommandsGateway) TodoDeleter {
	return &todoDeleter{todoCommandsGateway: todoCommandsGateway}
}

func (s *todoDeleter) Delete(ctx context.Context, in *input.TodoDeleter) (*output.TodoDeleter, error) {
	if err := s.todoCommandsGateway.Delete(ctx, in.ID); err != nil {
		return nil, fmt.Errorf("todoDeleter.Delete: %w", err)
	}

	return &output.TodoDeleter{}, nil
}
