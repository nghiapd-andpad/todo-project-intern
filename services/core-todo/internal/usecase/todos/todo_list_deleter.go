package todos

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoListDeleter interface {
	Delete(ctx context.Context, in *input.TodoListDeleter) (*output.TodoListDeleter, error)
}

type todoListDeleter struct {
	todoListCommandsGateway gateway.TodoListCommandsGateway
}

func NewTodoListDeleter(todoListCommandsGateway gateway.TodoListCommandsGateway) TodoListDeleter {
	return &todoListDeleter{todoListCommandsGateway: todoListCommandsGateway}
}

func (s *todoListDeleter) Delete(ctx context.Context, in *input.TodoListDeleter) (*output.TodoListDeleter, error) {
	if err := s.todoListCommandsGateway.Delete(ctx, in.ID); err != nil {
		return nil, fmt.Errorf("todoListDeleter.Delete: %w", err)
	}

	return &output.TodoListDeleter{}, nil
}
