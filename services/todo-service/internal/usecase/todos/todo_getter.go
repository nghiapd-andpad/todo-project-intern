package todos

import (
	"context"
	"errors"

	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/usecase/todos/output"
)

type TodoGetter interface {
	Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error)
}

// implementation of TodoGetter
type todoGetter struct {
	todoQueriesGateway gateway.TodoQueriesGateway
}

// NewTodoGetter creates a new instance of TodoGetter
func NewTodoGetter(todoQueriesGateway gateway.TodoQueriesGateway) TodoGetter {
	return &todoGetter{
		todoQueriesGateway: todoQueriesGateway,
	}
}

func (s *todoGetter) Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error) {
	if in == nil {
		return nil, errors.New("input is nil")
	}

	todo, err := s.todoQueriesGateway.Get(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if todo == nil {
		return nil, errors.New("todo not found")
	}

	return &output.TodoGetter{
		Todo: &entity.Todo{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      todo.Status,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		},
	}, nil
}
