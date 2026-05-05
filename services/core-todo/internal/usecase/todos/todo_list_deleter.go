package todos

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoListDeleter struct {
	todoListCommandsGateway gateway.TodoListCommandsGateway
}

func NewTodoListDeleter(todoListCommandsGateway gateway.TodoListCommandsGateway) *TodoListDeleter {
	return &TodoListDeleter{todoListCommandsGateway: todoListCommandsGateway}
}

func (s *TodoListDeleter) Delete(ctx context.Context, in *input.TodoListDeleter) (*output.TodoListDeleter, error) {
	if err := s.todoListCommandsGateway.Delete(ctx, in.ID); err != nil {
		return nil, fmt.Errorf("TodoListDeleter.Delete: %w", err)
	}

	return &output.TodoListDeleter{}, nil
}
