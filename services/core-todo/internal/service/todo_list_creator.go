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
	todoList := &entity.TodoList{
		Name:    in.Name,
		OwnerID: in.RequesterID,
	}

	var created *entity.TodoList

	err := s.transactor.Transaction(ctx, func(txCtx context.Context) error {
		// Find existing idempotency record by technical key.
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
			switch record.Status {
			case gatewayoutput.IdempotencyStatusCompleted:
				if record.ResourceID == nil {
					return entity.NewConflict("idempotency record completed without resource id")
				}

				// Replay by loading created resource.
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

			case gatewayoutput.IdempotencyStatusProcessing:
				return entity.NewConflict("idempotency request is still processing")

			default:
				return entity.NewConflict("idempotency request is in invalid state")
			}
		}

		// First execution: create PROCESSING record.
		record, err = s.idempotencyGateway.CreateProcessing(txCtx, &gatewayinput.CreateIdempotencyRecord{
			UserID:         in.RequesterID,
			Operation:      createTodoListOperation,
			IdempotencyKey: *in.IdempotencyKey,
			ExpiresAt:      time.Now().UTC().Add(24 * time.Hour),
		})
		if err != nil {
			return err
		}

		// Create resource.
		created, err = s.todoListCommandsGateway.Create(txCtx, todoList)
		if err != nil {
			return err
		}

		// Mark idempotency record as completed.
		if err := s.idempotencyGateway.MarkCompleted(
			txCtx,
			record.ID,
			"todo_list",
			int64(created.ID),
		); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("TodoListCreator.Create: %w", err)
	}

	return &output.TodoListCreator{TodoList: created}, nil
}
