package todos

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoListLister interface {
	List(ctx context.Context, in *input.TodoListLister) (*output.TodoListLister, error)
}

type todoListLister struct {
	todoListQueriesGateway gateway.TodoListQueriesGateway
}

func NewTodoListLister(todoListQueriesGateway gateway.TodoListQueriesGateway) TodoListLister {
	return &todoListLister{todoListQueriesGateway: todoListQueriesGateway}
}

func (s *todoListLister) List(ctx context.Context, in *input.TodoListLister) (*output.TodoListLister, error) {
	todoLists, total, err := s.todoListQueriesGateway.List(ctx, in.Opts)
	if err != nil {
		return nil, fmt.Errorf("todoListLister.List: %w", err)
	}

	return &output.TodoListLister{
		TodoLists: todoLists,
		Total:     total,
	}, nil
}
