package service

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

type TodoListCreator struct {
	todoListCommandsGateway gateway.TodoListCommandsGateway
}

func NewTodoListCreator(todoListCommandsGateway gateway.TodoListCommandsGateway) *TodoListCreator {
	return &TodoListCreator{todoListCommandsGateway: todoListCommandsGateway}
}

func (s *TodoListCreator) Create(ctx context.Context, in *input.TodoListCreator) (*output.TodoListCreator, error) {
	todoList := &entity.TodoList{
		Name:    in.Name,
		OwnerID: in.OwnerID,
	}

	created, err := s.todoListCommandsGateway.Create(ctx, todoList)
	if err != nil {
		return nil, fmt.Errorf("TodoListCreator.Create: %w", err)
	}

	return &output.TodoListCreator{TodoList: created}, nil
}
