// Package service contains business logic implementations for todo use cases.
package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

type TodoCreator struct {
	todoListQueriesGateway gateway.TodoListQueriesGateway
	todoCommandsGateway    gateway.TodoCommandsGateway
	cfg                    *config.Config
}

func NewTodoCreator(todoListQueriesGateway gateway.TodoListQueriesGateway, todoCommandsGateway gateway.TodoCommandsGateway, cfg *config.Config) *TodoCreator {
	return &TodoCreator{
		todoListQueriesGateway: todoListQueriesGateway,
		todoCommandsGateway:    todoCommandsGateway,
		cfg:                    cfg,
	}
}

func (s *TodoCreator) Create(ctx context.Context, in *input.TodoCreator) (*output.TodoCreator, error) {
	// Add feature balcklist
	if s.cfg.TodoBlacklistEnabled {
		if err := s.checkBlacklist(in.Title); err != nil {
			return nil, err
		}
	}

	// Find todo list
	todoList, err := s.todoListQueriesGateway.Get(ctx, in.TodoListID)
	if err != nil {
		return nil, fmt.Errorf("TodoCreator.Create: %w", err)
	}
	if todoList == nil {
		return nil, entity.NewNotFound("todo list not found")
	}

	if todoList.OwnerID != in.RequesterID {
		return nil, entity.NewAuthZ("you do not have permission to create todo in this list").
			WithDetail("owner_id", fmt.Sprintf("%d", todoList.OwnerID)).
			WithDetail("requester_id", fmt.Sprintf("%d", in.RequesterID))
	}

	todo := &entity.Todo{
		TodoListID:  in.TodoListID,
		Title:       in.Title,
		Description: in.Description,
		Status:      entity.TodoStatusPending,
		Priority:    in.Priority,
		DueDate:     in.DueDate,
		AssigneeID:  in.AssigneeID,
	}

	created, err := s.todoCommandsGateway.Create(ctx, todo)
	if err != nil {
		return nil, fmt.Errorf("TodoCreator.Create: %w", err)
	}

	return &output.TodoCreator{Todo: created}, nil
}

func (s *TodoCreator) checkBlacklist(title string) error {
	titleLower := strings.ToLower(title)
	for _, blocked := range s.cfg.TodoTitleBlacklist {
		if strings.Contains(titleLower, strings.ToLower(blocked)) {
			return entity.NewInvalidParameter("todo title contains a blacklisted word").
				WithDetail("title", title)
		}
	}
	return nil
}
