package todos

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoListCreator interface {
	Create(ctx context.Context, in *input.TodoListCreator) (*output.TodoListCreator, error)
}

type todoListCreator struct {
	todoListCommandsGateway gateway.TodoListCommandsGateway
}

func NewTodoListCreator(todoListCommandsGateway gateway.TodoListCommandsGateway) TodoListCreator {
	return &todoListCreator{todoListCommandsGateway: todoListCommandsGateway}
}

func (s *todoListCreator) Create(ctx context.Context, in *input.TodoListCreator) (*output.TodoListCreator, error) {
	todoList := &entity.TodoList{
		Name:    in.Name,
		OwnerID: in.OwnerID,
	}

	created, err := s.todoListCommandsGateway.Create(ctx, todoList)
	if err != nil {
		return nil, fmt.Errorf("todoListCreator.Create: %w", err)
	}

	return &output.TodoListCreator{TodoList: created}, nil
}
