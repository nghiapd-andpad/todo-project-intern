package todos

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoListUpdater interface {
	Update(ctx context.Context, in *input.TodoListUpdater) (*output.TodoListUpdater, error)
}

type todoListUpdater struct {
	todoListCommandsGateway gateway.TodoListCommandsGateway
	todoListQueriesGateway  gateway.TodoListQueriesGateway
}

func NewTodoListUpdater(
	todoListCommandsGateway gateway.TodoListCommandsGateway,
	todoListQueriesGateway gateway.TodoListQueriesGateway,
) TodoListUpdater {
	return &todoListUpdater{
		todoListCommandsGateway: todoListCommandsGateway,
		todoListQueriesGateway:  todoListQueriesGateway,
	}
}

func (s *todoListUpdater) Update(ctx context.Context, in *input.TodoListUpdater) (*output.TodoListUpdater, error) {
	todoList, err := s.todoListQueriesGateway.Get(ctx, in.ID)
	if err != nil {
		return nil, fmt.Errorf("todoListUpdater.Get: %w", err)
	}
	if todoList == nil {
		return nil, entity.NewNotFound("todo list not found").
			WithDetail("todo_list_id", fmt.Sprintf("%d", in.ID))
	}

	if in.Name != nil {
		todoList.Name = *in.Name
	}

	updated, err := s.todoListCommandsGateway.Update(ctx, todoList)
	if err != nil {
		return nil, fmt.Errorf("todoListUpdater.Update: %w", err)
	}

	return &output.TodoListUpdater{TodoList: updated}, nil
}
