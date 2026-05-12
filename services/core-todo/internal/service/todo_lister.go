package service

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/shared/pagination"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

type TodoLister struct {
	todoListQueriesGateway gateway.TodoListQueriesGateway
	todoQueriesGateway     gateway.TodoQueriesGateway
}

func NewTodoLister(todoListQueriesGateway gateway.TodoListQueriesGateway, todoQueriesGateway gateway.TodoQueriesGateway) *TodoLister {
	return &TodoLister{
		todoListQueriesGateway: todoListQueriesGateway,
		todoQueriesGateway:     todoQueriesGateway,
	}
}

func (s *TodoLister) List(ctx context.Context, in *input.TodoLister) (*output.TodoLister, error) {
	todoList, err := s.todoListQueriesGateway.Get(ctx, in.TodoListID)
	if err != nil {
		return nil, fmt.Errorf("TodoLister.List: %w", err)
	}
	if todoList == nil {
		return nil, entity.NewNotFound("todo list not found")
	}

	opts := s.buildGatewayOpts(in, todoList)

	items, total, err := s.todoQueriesGateway.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("TodoLister.List: %w", err)
	}

	return &output.TodoLister{
		Page: pagination.New(items, total, in.Offset, in.Limit),
	}, nil
}

func (s *TodoLister) buildGatewayOpts(in *input.TodoLister, todoList *entity.TodoList) *gatewayinput.ListTodosOptions {
	opts := &gatewayinput.ListTodosOptions{
		TodoListID:  in.TodoListID,
		Status:      in.Status,
		Priority:    in.Priority,
		TitleSearch: in.TitleSearch,
		Offset:      in.Offset,
		Limit:       in.Limit,
	}

	if todoList.OwnerID == in.RequesterID {
		opts.AssigneeOnly = nil
	} else {
		opts.AssigneeOnly = &in.RequesterID
	}

	return opts
}
