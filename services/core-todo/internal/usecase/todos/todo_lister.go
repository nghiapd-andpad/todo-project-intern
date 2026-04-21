package todos

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoLister interface {
	List(ctx context.Context, in *input.TodoLister) (*output.TodoLister, error)
}

type todoLister struct {
	todoQueriesGateway gateway.TodoQueriesGateway
}

func NewTodoLister(todoQueriesGateway gateway.TodoQueriesGateway) TodoLister {
	return &todoLister{todoQueriesGateway: todoQueriesGateway}
}

func (s *todoLister) List(ctx context.Context, in *input.TodoLister) (*output.TodoLister, error) {
	todos, total, err := s.todoQueriesGateway.List(ctx, in.Opts)
	if err != nil {
		return nil, fmt.Errorf("todoLister.List: %w", err)
	}

	return &output.TodoLister{
		Todos: todos,
		Total: total,
	}, nil
}
