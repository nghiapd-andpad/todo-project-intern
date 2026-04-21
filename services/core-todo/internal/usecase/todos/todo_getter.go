package todos

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoGetter interface {
	Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error)
}

type todoGetter struct {
	todoQueriesGateway gateway.TodoQueriesGateway
}

func NewTodoGetter(todoQueriesGateway gateway.TodoQueriesGateway) TodoGetter {
	return &todoGetter{todoQueriesGateway: todoQueriesGateway}
}

func (s *todoGetter) Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error) {
	todo, err := s.todoQueriesGateway.Get(ctx, in.ID)
	if err != nil {
		return nil, fmt.Errorf("todoGetter.Get: %w", err)
	}

	if todo == nil {
		return nil, entity.NewNotFound("todo not found").
			WithDetail("todo_id", fmt.Sprintf("%d", in.ID))
	}

	return &output.TodoGetter{Todo: todo}, nil
}
