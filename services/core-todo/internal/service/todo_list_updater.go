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
	transactor              gateway.Transactor
	todoListCommandsGateway gateway.TodoListCommandsGateway
	todoListQueriesGateway  gateway.TodoListQueriesGateway
}

func NewTodoListUpdater(
	transactor gateway.Transactor,
	todoListCommandsGateway gateway.TodoListCommandsGateway,
	todoListQueriesGateway gateway.TodoListQueriesGateway,
) *TodoListUpdater {
	return &TodoListUpdater{
		transactor:              transactor,
		todoListCommandsGateway: todoListCommandsGateway,
		todoListQueriesGateway:  todoListQueriesGateway,
	}
}

func (s *TodoListUpdater) Update(ctx context.Context, in *input.TodoListUpdater) (*output.TodoListUpdater, error) {
	var updated *entity.TodoList

	err := s.transactor.Transaction(ctx, func(txCtx context.Context) error {
		// Lock todo list row before applying update logic.
		todoList, err := s.todoListQueriesGateway.GetForUpdate(txCtx, in.TodoListID)
		if err != nil {
			return fmt.Errorf("TodoListUpdater.GetForUpdate: %w", err)
		}
		if todoList == nil {
			return entity.NewNotFound("todo list not found").
				WithDetail("todo_list_id", fmt.Sprintf("%d", in.TodoListID))
		}

		if todoList.OwnerID != in.RequesterID {
			return entity.NewAuthZ("you do not have permission to update this todo list").
				WithDetail("owner_id", fmt.Sprintf("%d", todoList.OwnerID)).
				WithDetail("requester_id", fmt.Sprintf("%d", in.RequesterID))
		}

		if in.Name != nil {
			todoList.Name = *in.Name
		}

		updated, err = s.todoListCommandsGateway.Update(txCtx, todoList)
		if err != nil {
			return fmt.Errorf("TodoListUpdater.Update: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("TodoListUpdater.Update: %w", err)
	}

	return &output.TodoListUpdater{TodoList: updated}, nil
}
