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
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/utils/idempotency"
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

	requestHash, err := idempotency.Hash(struct {
		TodoListID  entity.TodoListID
		RequesterID entity.UserID
		Title       string
		Description *string
		Priority    entity.Priority
		DueDate     *time.Time
		AssigneeID  *entity.UserID
	}{
		TodoListID:  in.TodoListID,
		RequesterID: in.RequesterID,
		Title:       in.Title,
		Description: in.Description,
		Priority:    in.Priority,
		DueDate:     in.DueDate,
		AssigneeID:  in.AssigneeID,
	})
	if err != nil {
		return nil, fmt.Errorf("TodoCreator.Create: build request hash: %w", err)
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

	// No idempotency key means normal create flow.
	if in.IdempotencyKey == nil || *in.IdempotencyKey == "" {
		created, err := s.todoCommandsGateway.Create(ctx, todo)
		if err != nil {
			return nil, fmt.Errorf("TodoCreator.Create: %w", err)
		}

		return &output.TodoCreator{Todo: created}, nil
	}

	var created *entity.Todo

	err = s.transactor.Transaction(ctx, func(txCtx context.Context) error {
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
			if record.RequestHash != requestHash {
				return entity.NewInvalidParameter("idempotency key reused with different request")
			}

			if record.Status == gatewayoutput.IdempotencyStatusCompleted {
				if record.ResourceID == nil {
					return entity.NewConflict("idempotency record completed without resource id")
				}

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
			}

			if record.Status == gatewayoutput.IdempotencyStatusProcessing {
				return entity.NewConflict("idempotency request is still processing")
			}

			return entity.NewConflict("idempotency request is in invalid state")
		}

		record, err = s.idempotencyGateway.CreateProcessing(txCtx, &gatewayinput.CreateIdempotencyRecord{
			UserID:         in.RequesterID,
			Operation:      createTodoOperation,
			IdempotencyKey: *in.IdempotencyKey,
			RequestHash:    requestHash,
			ExpiresAt:      time.Now().UTC().Add(24 * time.Hour),
		})
		if err != nil {
			return err
		}

		created, err = s.todoCommandsGateway.Create(txCtx, todo)
		if err != nil {
			return err
		}

		if err := s.idempotencyGateway.MarkCompleted(txCtx, record.ID, int64(created.ID)); err != nil {
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
