package todos

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoLister struct {
	todoQueriesGateway gateway.TodoQueriesGateway
}

func NewTodoLister(todoQueriesGateway gateway.TodoQueriesGateway) *TodoLister {
	return &TodoLister{todoQueriesGateway: todoQueriesGateway}
}

func (s *TodoLister) List(ctx context.Context, in *input.TodoLister) (*output.TodoLister, error) {
	in.Opts.SetDefaults()
	if err := in.Opts.Validate(); err != nil {
		return nil, err
	}

	todos, total, err := s.todoQueriesGateway.List(ctx, in.Opts)
	if err != nil {
		return nil, fmt.Errorf("TodoLister.List: %w", err)
	}

	return &output.TodoLister{
		Todos: todos,
		Total: total,
	}, nil
}
