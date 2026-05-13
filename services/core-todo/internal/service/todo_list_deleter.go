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
	transactor              gateway.Transactor
	todoListQueriesGateway  gateway.TodoListQueriesGateway
	todoListCommandsGateway gateway.TodoListCommandsGateway
	todoCommandsGateway     gateway.TodoCommandsGateway
}

func NewTodoListDeleter(
	transactor gateway.Transactor,
	todoListQueriesGateway gateway.TodoListQueriesGateway,
	todoListCommandsGateway gateway.TodoListCommandsGateway,
	todoCommandsGateway gateway.TodoCommandsGateway,
) *TodoListDeleter {
	return &TodoListDeleter{
		transactor:              transactor,
		todoListQueriesGateway:  todoListQueriesGateway,
		todoListCommandsGateway: todoListCommandsGateway,
		todoCommandsGateway:     todoCommandsGateway,
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

	err = s.transactor.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.todoCommandsGateway.DeleteByTodoListID(txCtx, in.TodoListID); err != nil {
			return fmt.Errorf("delete todos: %w", err)
		}

		if err := s.todoListCommandsGateway.Delete(txCtx, in.TodoListID); err != nil {
			return fmt.Errorf("delete todo list: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("TodoListDeleter.Delete: %w", err)
	}

	return &output.TodoListDeleter{}, nil
}
