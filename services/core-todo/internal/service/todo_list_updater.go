package service

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

type TodoListUpdater struct {
	todoListCommandsGateway gateway.TodoListCommandsGateway
	todoListQueriesGateway  gateway.TodoListQueriesGateway
}

func NewTodoListUpdater(
	todoListCommandsGateway gateway.TodoListCommandsGateway,
	todoListQueriesGateway gateway.TodoListQueriesGateway,
) *TodoListUpdater {
	return &TodoListUpdater{
		todoListCommandsGateway: todoListCommandsGateway,
		todoListQueriesGateway:  todoListQueriesGateway,
	}
}

func (s *TodoListUpdater) Update(ctx context.Context, in *input.TodoListUpdater) (*output.TodoListUpdater, error) {
	todoList, err := s.todoListQueriesGateway.Get(ctx, in.TodoListID)
	if err != nil {
		return nil, fmt.Errorf("TodoListUpdater.Get: %w", err)
	}
	if todoList == nil {
		return nil, entity.NewNotFound("todo list not found").
			WithDetail("todo_list_id", fmt.Sprintf("%d", in.TodoListID))
	}

	if todoList.OwnerID != in.RequesterID {
		return nil, entity.NewAuthZ("you do not have permission to update this todo list").
			WithDetail("owner_id", fmt.Sprintf("%d", todoList.OwnerID)).
			WithDetail("requester_id", fmt.Sprintf("%d", in.RequesterID))
	}

	todoList.Version = in.Version

	if in.Name != nil {
		todoList.Name = *in.Name
	}

	updated, err := s.todoListCommandsGateway.Update(ctx, todoList)
	if err != nil {
		return nil, fmt.Errorf("TodoListUpdater.Update: %w", err)
	}

	return &output.TodoListUpdater{TodoList: updated}, nil
}
