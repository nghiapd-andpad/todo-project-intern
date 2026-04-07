package todos

import (
	"context"

	"github.com/nghiaphunng18/todos/internal/gateway"
	"github.com/nghiaphunng18/todos/internal/usecase/todos/input"
	"github.com/nghiaphunng18/todos/internal/usecase/todos/output"
)

type TodoListReader interface {
	List(ctx context.Context, in *input.TodoListReader) (*output.TodoListReader, error)
}

type todoListReader struct {
	todoQueriesGateway gateway.TodoQueriesGateway
}

func NewTodoListReader(todoQueriesGateway gateway.TodoQueriesGateway) TodoListReader {
	return &todoListReader{
		todoQueriesGateway: todoQueriesGateway,
	}
}

func (s *todoListReader) List(ctx context.Context, in *input.TodoListReader) (*output.TodoListReader, error) {
	todos, err := s.todoQueriesGateway.List(ctx)
	if err != nil {
		return nil, err
	}

	return &output.TodoListReader{
		Todos: todos,
	}, nil
}
