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

type TodoCreator interface {
	Create(ctx context.Context, in *input.TodoCreator) (*output.TodoCreator, error)
}

type todoCreator struct {
	todoCommandsGateway gateway.TodoCommandsGateway
}

func NewTodoCreator(todoCommandsGateway gateway.TodoCommandsGateway) TodoCreator {
	return &todoCreator{todoCommandsGateway: todoCommandsGateway}
}

func (s *todoCreator) Create(ctx context.Context, in *input.TodoCreator) (*output.TodoCreator, error) {
	var dueDate *time.Time
	if in.DueDate != nil {
		parsed, err := time.Parse("2006-01-02", *in.DueDate)
		if err != nil {
			return nil, entity.NewInvalidParameter("invalid due_date format, expected YYYY-MM-DD").
				WithDetail("due_date", *in.DueDate)
		}
		dueDate = &parsed
	}

	todo := &entity.Todo{
		TodoListID:  in.TodoListID,
		Title:       in.Title,
		Description: in.Description,
		Status:      entity.TodoStatusPending,
		Priority:    in.Priority,
		DueDate:     dueDate,
		AssigneeID:  in.AssigneeID,
		CreatorID:   in.CreatorID,
	}

	created, err := s.todoCommandsGateway.Create(ctx, todo)
	if err != nil {
		return nil, fmt.Errorf("todoCreator.Create: %w", err)
	}

	return &output.TodoCreator{Todo: created}, nil
}
