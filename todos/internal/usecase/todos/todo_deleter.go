package todos

import (
	"context"

	"github.com/nghiaphunng18/todos/internal/gateway"
	"github.com/nghiaphunng18/todos/internal/usecase/todos/input"
	"github.com/nghiaphunng18/todos/internal/usecase/todos/output"
)

type TodoDeleter interface {
	Delete(ctx context.Context, in *input.TodoDeleter) (*output.TodoDeleter, error)
}

type todoDeleter struct {
	todoCommandsGateway gateway.TodoCommandsGateway
}

func NewTodoDeleter(todoCommandsGateway gateway.TodoCommandsGateway) TodoDeleter {
	return &todoDeleter{
		todoCommandsGateway: todoCommandsGateway,
	}
}

func (s *todoDeleter) Delete(ctx context.Context, in *input.TodoDeleter) (*output.TodoDeleter, error) {
	err := s.todoCommandsGateway.Delete(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	return &output.TodoDeleter{
		Success: true,
	}, nil
}
