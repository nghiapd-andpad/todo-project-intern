package service

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

type TodoListDeleter struct {
	todoListQueriesGateway  gateway.TodoListQueriesGateway
	todoListCommandsGateway gateway.TodoListCommandsGateway
}

func NewTodoListDeleter(todoListQueriesGateway gateway.TodoListQueriesGateway, todoListCommandsGateway gateway.TodoListCommandsGateway) *TodoListDeleter {
	return &TodoListDeleter{
		todoListQueriesGateway:  todoListQueriesGateway,
		todoListCommandsGateway: todoListCommandsGateway,
	}
}

func (s *TodoListDeleter) Delete(ctx context.Context, in *input.TodoListDeleter) (*output.TodoListDeleter, error) {
	todoList, err := s.todoListQueriesGateway.Get(ctx, in.TodoListID)
	if err != nil {
		return nil, fmt.Errorf("TodoListDeleter.Delete: %w", err)
	}
	if todoList == nil {
		return nil, entity.NewNotFound("todo list not found").
			WithDetail("todo_list_id", fmt.Sprintf("%d", in.TodoListID))
	}

	if todoList.OwnerID != in.RequesterID {
		return nil, entity.NewAuthZ("you do not have permission to delete this todo list").
			WithDetail("owner_id", fmt.Sprintf("%d", todoList.OwnerID)).
			WithDetail("requester_id", fmt.Sprintf("%d", in.RequesterID))
	}

	if err := s.todoListCommandsGateway.Delete(ctx, in.TodoListID); err != nil {
		return nil, fmt.Errorf("TodoListDeleter.Delete: %w", err)
	}

	return &output.TodoListDeleter{}, nil
}
