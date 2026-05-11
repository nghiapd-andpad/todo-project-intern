package service

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

type TodoListLister struct {
	todoListQueriesGateway gateway.TodoListQueriesGateway
}

func NewTodoListLister(todoListQueriesGateway gateway.TodoListQueriesGateway) *TodoListLister {
	return &TodoListLister{todoListQueriesGateway: todoListQueriesGateway}
}

func (s *TodoListLister) List(ctx context.Context, in *input.TodoListLister) (*output.TodoListLister, error) {
	todoLists, total, err := s.todoListQueriesGateway.List(ctx, &in.Opts)
	if err != nil {
		return nil, fmt.Errorf("TodoListLister.List: %w", err)
	}

	return &output.TodoListLister{
		TodoLists: todoLists,
		Total:     total,
	}, nil
}
