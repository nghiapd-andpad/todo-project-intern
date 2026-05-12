package service

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

type TodoGetter struct {
	todoListQueriesGateway gateway.TodoListQueriesGateway
	todoQueriesGateway     gateway.TodoQueriesGateway
}

func NewTodoGetter(todoListQueriesGateway gateway.TodoListQueriesGateway, todoQueriesGateway gateway.TodoQueriesGateway) *TodoGetter {
	return &TodoGetter{
		todoListQueriesGateway: todoListQueriesGateway,
		todoQueriesGateway:     todoQueriesGateway,
	}
}

func (s *TodoGetter) Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error) {
	todo, err := s.todoQueriesGateway.Get(ctx, in.TodoID, in.TodoListID)
	if err != nil {
		return nil, fmt.Errorf("TodoGetter.Get: %w", err)
	}
	if todo == nil {
		return nil, entity.NewNotFound("todo not found")
	}

	todoList, err := s.todoListQueriesGateway.Get(ctx, in.TodoListID)
	if err != nil {
		return nil, fmt.Errorf("TodoGetter.Get: %w", err)
	}
	if todoList == nil {
		return nil, entity.NewNotFound("todo list not found")
	}

	if todoList.OwnerID != in.RequesterID && (todo.AssigneeID == nil || *todo.AssigneeID != in.RequesterID) {
		return nil, entity.NewAuthZ("you do not have permission to view this todo")
	}

	return &output.TodoGetter{Todo: todo}, nil
}
