package todos

import (
	"context"

	"github.com/nghiaphunng18/todos/internal/gateway"
	"github.com/nghiaphunng18/todos/internal/usecase/todos/input"
	"github.com/nghiaphunng18/todos/internal/usecase/todos/output"
)

type TodoLister interface {
	List(ctx context.Context, in *input.TodoLister) (*output.TodoLister, error)
}

type todoListReader struct {
	todoQueriesGateway gateway.TodoQueriesGateway
}

func NewTodoListReader(todoQueriesGateway gateway.TodoQueriesGateway) TodoLister {
	return &todoListReader{
		todoQueriesGateway: todoQueriesGateway,
	}
}

func (s *todoListReader) List(ctx context.Context, in *input.TodoLister) (*output.TodoLister, error) {
	todos, err := s.todoQueriesGateway.List(ctx)
	if err != nil {
		return nil, err
	}

	return &output.TodoLister{
		Todos: todos,
	}, nil
}
