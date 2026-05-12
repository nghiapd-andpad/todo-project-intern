package service

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

type TodoDeleter struct {
	todoListQueriesGateway gateway.TodoListQueriesGateway
	todoQueriesGateway     gateway.TodoQueriesGateway
	todoCommandsGateway    gateway.TodoCommandsGateway
}

func NewTodoDeleter(
	todoListQueriesGateway gateway.TodoListQueriesGateway,
	todoQueriesGateway gateway.TodoQueriesGateway,
	todoCommandsGateway gateway.TodoCommandsGateway,
) *TodoDeleter {
	return &TodoDeleter{
		todoListQueriesGateway: todoListQueriesGateway,
		todoQueriesGateway:     todoQueriesGateway,
		todoCommandsGateway:    todoCommandsGateway,
	}
}

func (s *TodoDeleter) Delete(ctx context.Context, in *input.TodoDeleter) (*output.TodoDeleter, error) {
	todo, err := s.todoQueriesGateway.Get(ctx, in.TodoID, in.TodoListID)
	if err != nil {
		return nil, fmt.Errorf("TodoDeleter.Delete: %w", err)
	}
	if todo == nil {
		return nil, entity.NewNotFound("todo not found")
	}

	todoList, err := s.todoListQueriesGateway.Get(ctx, in.TodoListID)
	if err != nil {
		return nil, fmt.Errorf("TodoDeleter.Delete: %w", err)
	}
	if todoList == nil {
		return nil, entity.NewNotFound("todo list not found")
	}

	if todoList.OwnerID != in.RequesterID {
		return nil, entity.NewAuthZ("you do not have permission to delete this todo")
	}

	if err := s.todoCommandsGateway.Delete(ctx, in.TodoID); err != nil {
		return nil, fmt.Errorf("TodoDeleter.Delete: %w", err)
	}

	return &output.TodoDeleter{}, nil
}
