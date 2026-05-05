package todos

import (
	"context"
	"fmt"
	"time"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoUpdater struct {
	todoCommandsGateway gateway.TodoCommandsGateway
	todoQueriesGateway  gateway.TodoQueriesGateway
}

func NewTodoUpdater(
	todoCommandsGateway gateway.TodoCommandsGateway,
	todoQueriesGateway gateway.TodoQueriesGateway,
) *TodoUpdater {
	return &TodoUpdater{
		todoCommandsGateway: todoCommandsGateway,
		todoQueriesGateway:  todoQueriesGateway,
	}
}

func (s *TodoUpdater) Update(ctx context.Context, in *input.TodoUpdater) (*output.TodoUpdater, error) {
	todo, err := s.todoQueriesGateway.Get(ctx, in.ID)
	if err != nil {
		return nil, fmt.Errorf("TodoUpdater.Get: %w", err)
	}
	if todo == nil {
		return nil, entity.NewNotFound("todo not found").
			WithDetail("todo_id", fmt.Sprintf("%d", in.ID))
	}

	if in.Title != nil {
		todo.Title = *in.Title
	}
	if in.Description != nil {
		todo.Description = in.Description
	}
	if in.Status != nil {
		todo.Status = *in.Status
	}
	if in.Priority != nil {
		todo.Priority = *in.Priority
	}
	if in.AssigneeID != nil {
		todo.AssigneeID = in.AssigneeID
	}
	if in.DueDate != nil {
		parsed, err := time.Parse("2006-01-02", *in.DueDate)
		if err != nil {
			return nil, entity.NewInvalidParameter("invalid due_date format, expected YYYY-MM-DD").
				WithDetail("due_date", *in.DueDate)
		}
		todo.DueDate = &parsed
	}

	updated, err := s.todoCommandsGateway.Update(ctx, todo)
	if err != nil {
		return nil, fmt.Errorf("TodoUpdater.Update: %w", err)
	}

	return &output.TodoUpdater{Todo: updated}, nil
}
