package service

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/shared/pagination"
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
	opts := s.buildGatewayOpts(in)

	todoLists, total, err := s.todoListQueriesGateway.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("TodoListLister.List: %w", err)
	}

	return &output.TodoListLister{
		Page: pagination.New(todoLists, total, in.Offset, in.Limit),
	}, nil
}

func (s *TodoListLister) buildGatewayOpts(in *input.TodoListLister) *gatewayinput.ListTodoListsOptions {
	opts := &gatewayinput.ListTodoListsOptions{
		NameSearch: in.NameSearch,
		Offset:     in.Offset,
		Limit:      in.Limit,
	}

	switch in.Filter {
	case input.TodoListFilterOwned:
		opts.OwnerID = &in.RequesterID

	case input.TodoListFilterAssigned:
		opts.AssigneeID = &in.RequesterID

	default:
		opts.OwnerID = &in.RequesterID
		opts.AssigneeID = &in.RequesterID
	}

	return opts
}
