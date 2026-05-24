package service

import (
	"context"
	"fmt"
	"time"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	gatewayoutput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/output"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/utils/idempotency"
)

const createTodoListOperation = "CREATE_TODO_LIST"

type TodoListCreator struct {
	transactor              gateway.Transactor
	idempotencyGateway      gateway.IdempotencyGateway
	todoListQueriesGateway  gateway.TodoListQueriesGateway
	todoListCommandsGateway gateway.TodoListCommandsGateway
}

func NewTodoListCreator(
	transactor gateway.Transactor,
	idempotencyGateway gateway.IdempotencyGateway,
	todoListQueriesGateway gateway.TodoListQueriesGateway,
	todoListCommandsGateway gateway.TodoListCommandsGateway,
) *TodoListCreator {
	return &TodoListCreator{
		transactor:              transactor,
		idempotencyGateway:      idempotencyGateway,
		todoListQueriesGateway:  todoListQueriesGateway,
		todoListCommandsGateway: todoListCommandsGateway,
	}
}

func (s *TodoListCreator) Create(ctx context.Context, in *input.TodoListCreator) (*output.TodoListCreator, error) {
	requestHash, err := idempotency.Hash(struct {
		RequesterID entity.UserID
		Name        string
	}{
		RequesterID: in.RequesterID,
		Name:        in.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("TodoListCreator.Create: build request hash: %w", err)
	}

	todoList := &entity.TodoList{
		Name:    in.Name,
		OwnerID: in.RequesterID,
	}

	if in.IdempotencyKey == nil || *in.IdempotencyKey == "" {
		created, err := s.todoListCommandsGateway.Create(ctx, todoList)
		if err != nil {
			return nil, fmt.Errorf("TodoListCreator.Create: %w", err)
		}

		return &output.TodoListCreator{TodoList: created}, nil
	}

	var created *entity.TodoList

	err = s.transactor.Transaction(ctx, func(txCtx context.Context) error {
		record, err := s.idempotencyGateway.Find(
			txCtx,
			in.RequesterID,
			createTodoListOperation,
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

				replayed, err := s.todoListQueriesGateway.Get(
					txCtx,
					entity.TodoListID(*record.ResourceID),
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
			Operation:      createTodoListOperation,
			IdempotencyKey: *in.IdempotencyKey,
			RequestHash:    requestHash,
			ExpiresAt:      time.Now().UTC().Add(24 * time.Hour),
		})
		if err != nil {
			return err
		}

		created, err = s.todoListCommandsGateway.Create(txCtx, todoList)
		if err != nil {
			return err
		}

		if err := s.idempotencyGateway.MarkCompleted(txCtx, record.ID, int64(created.ID)); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("TodoListCreator.Create: %w", err)
	}

	return &output.TodoListCreator{TodoList: created}, nil
}
