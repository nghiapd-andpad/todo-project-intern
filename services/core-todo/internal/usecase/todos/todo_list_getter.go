package todos

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoListGetter interface {
	Get(ctx context.Context, in *input.TodoListGetter) (*output.TodoListGetter, error)
}

type todoListGetter struct {
	todoListQueriesGateway gateway.TodoListQueriesGateway
}

func NewTodoListGetter(todoListQueriesGateway gateway.TodoListQueriesGateway) TodoListGetter {
	return &todoListGetter{todoListQueriesGateway: todoListQueriesGateway}
}

func (s *todoListGetter) Get(ctx context.Context, in *input.TodoListGetter) (*output.TodoListGetter, error) {
	todoList, err := s.todoListQueriesGateway.Get(ctx, in.ID)
	if err != nil {
		return nil, fmt.Errorf("todoListGetter.Get: %w", err)
	}

	if todoList == nil {
		return nil, entity.NewNotFound("todo list not found").
			WithDetail("todo_list_id", fmt.Sprintf("%d", in.ID))
	}

	return &output.TodoListGetter{TodoList: todoList}, nil
}
