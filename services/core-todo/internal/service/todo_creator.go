// Package service contains business logic implementations for todo use cases.
package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	gatewayoutput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/output"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
)

const createTodoOperation = "CREATE_TODO"

type TodoCreator struct {
	cfg                    *config.Config
	transactor             gateway.Transactor
	idempotencyGateway     gateway.IdempotencyGateway
	todoQueriesGateway     gateway.TodoQueriesGateway
	todoCommandsGateway    gateway.TodoCommandsGateway
	todoListQueriesGateway gateway.TodoListQueriesGateway
}

func NewTodoCreator(cfg *config.Config, transactor gateway.Transactor, idempotencyGateway gateway.IdempotencyGateway, todoQueriesGateway gateway.TodoQueriesGateway, todoCommandsGateway gateway.TodoCommandsGateway, todoListQueriesGateway gateway.TodoListQueriesGateway) *TodoCreator {
	return &TodoCreator{
		cfg:                    cfg,
		transactor:             transactor,
		idempotencyGateway:     idempotencyGateway,
		todoQueriesGateway:     todoQueriesGateway,
		todoCommandsGateway:    todoCommandsGateway,
		todoListQueriesGateway: todoListQueriesGateway,
	}
}

func (s *TodoCreator) Create(ctx context.Context, in *input.TodoCreator) (*output.TodoCreator, error) {
	if s.cfg.TodoBlacklistEnabled {
		if err := s.checkBlacklist(in.Title); err != nil {
			return nil, err
		}
	}

	// Check permissstion first.
	todoList, err := s.todoListQueriesGateway.Get(ctx, in.TodoListID)
	if err != nil {
		return nil, fmt.Errorf("TodoCreator.Create: %w", err)
	}
	if todoList == nil {
		return nil, entity.NewNotFound("todo list not found")
	}
	if todoList.OwnerID != in.RequesterID {
		return nil, entity.NewAuthZ("you do not have permission to create todo in this list")
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

	var created *entity.Todo

	err = s.transactor.Transaction(ctx, func(txCtx context.Context) error {
		// Find existing idempotency record by technical key.
		record, err := s.idempotencyGateway.Find(
			txCtx,
			in.RequesterID,
			createTodoOperation,
			*in.IdempotencyKey,
		)
		if err != nil {
			return err
		}

		if record != nil {
			switch record.Status {
			case gatewayoutput.IdempotencyStatusCompleted:
				if record.ResourceID == nil {
					return entity.NewConflict("idempotency record completed without resource id")
				}

				// Replay by loading created resource.
				replayed, err := s.todoQueriesGateway.Get(
					txCtx,
					entity.TodoID(*record.ResourceID),
					in.TodoListID,
				)
				if err != nil {
					return err
				}
				if replayed == nil {
					return entity.NewNotFound("idempotency resource not found")
				}

				created = replayed
				return nil

			case gatewayoutput.IdempotencyStatusProcessing:
				return entity.NewConflict("idempotency request is still processing")

			default:
				return entity.NewConflict("idempotency request is in invalid state")
			}
		}

		// If record not found: create PROCESSING record.
		record, err = s.idempotencyGateway.CreateProcessing(txCtx, &gatewayinput.CreateIdempotencyRecord{
			UserID:         in.RequesterID,
			Operation:      createTodoOperation,
			IdempotencyKey: *in.IdempotencyKey,
			ExpiresAt:      time.Now().UTC().Add(24 * time.Hour),
		})
		if err != nil {
			return err
		}

		// Create resource.
		created, err = s.todoCommandsGateway.Create(txCtx, todo)
		if err != nil {
			return err
		}

		// Mark idempotency record as completed.
		if err := s.idempotencyGateway.MarkCompleted(
			txCtx,
			record.ID,
			"todo",
			int64(created.ID),
		); err != nil {
			return err
		}

		return nil
	})
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
