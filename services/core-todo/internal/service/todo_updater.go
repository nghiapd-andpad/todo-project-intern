package service

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

type TodoUpdater struct {
	transactor             gateway.Transactor
	todoListQueriesGateway gateway.TodoListQueriesGateway
	todoQueriesGateway     gateway.TodoQueriesGateway
	todoCommandsGateway    gateway.TodoCommandsGateway
}

func NewTodoUpdater(
	transactor gateway.Transactor,
	todoListQueriesGateway gateway.TodoListQueriesGateway,
	todoQueriesGateway gateway.TodoQueriesGateway,
	todoCommandsGateway gateway.TodoCommandsGateway,
) *TodoUpdater {
	return &TodoUpdater{
		transactor:             transactor,
		todoListQueriesGateway: todoListQueriesGateway,
		todoQueriesGateway:     todoQueriesGateway,
		todoCommandsGateway:    todoCommandsGateway,
	}
}

func (s *TodoUpdater) Update(ctx context.Context, in *input.TodoUpdater) (*output.TodoUpdater, error) {
	var updated *entity.Todo

	err := s.transactor.Transaction(ctx, func(txCtx context.Context) error {
		// Lock todo row before applying update logic.
		todo, err := s.todoQueriesGateway.GetForUpdate(txCtx, in.TodoID, in.TodoListID)
		if err != nil {
			return fmt.Errorf("TodoUpdater.GetForUpdate: %w", err)
		}
		if todo == nil {
			return entity.NewNotFound("todo not found")
		}

		todoList, err := s.todoListQueriesGateway.Get(txCtx, in.TodoListID)
		if err != nil {
			return fmt.Errorf("TodoUpdater.GetTodoList: %w", err)
		}
		if todoList == nil {
			return entity.NewNotFound("todo list not found")
		}

		isOwner := todoList.OwnerID == in.RequesterID
		isAssignee := todo.AssigneeID != nil && *todo.AssigneeID == in.RequesterID
		if !isOwner && !isAssignee {
			return entity.NewAuthZ("you do not have permission to update this todo")
		}

		if isOwner {
			applyAllFields(todo, in.Fields)
		} else {
			// Assignee can only update status.
			applyAssigneeFields(todo, in.Fields)
		}

		updated, err = s.todoCommandsGateway.Update(txCtx, todo)
		if err != nil {
			return fmt.Errorf("TodoUpdater.Update: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("TodoUpdater.Update: %w", err)
	}

	return &output.TodoUpdater{Todo: updated}, nil
}

func applyAllFields(todo *entity.Todo, fields input.UpdateTodoFields) {
	if fields.Title != nil {
		todo.Title = *fields.Title
	}
	if fields.Description != nil {
		todo.Description = fields.Description
	}
	if fields.Status != nil {
		todo.Status = *fields.Status
	}
	if fields.Priority != nil {
		todo.Priority = *fields.Priority
	}
	if fields.DueDate != nil {
		todo.DueDate = fields.DueDate
	}
	if fields.AssigneeID != nil {
		todo.AssigneeID = fields.AssigneeID
	}
}

func applyAssigneeFields(todo *entity.Todo, fields input.UpdateTodoFields) {
	if fields.Status != nil {
		todo.Status = *fields.Status
	}
}
