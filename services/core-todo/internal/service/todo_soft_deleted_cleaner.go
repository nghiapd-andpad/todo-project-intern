package service

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/output"
	logutil "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/utils/logger"
)

type TodoSoftDeletedCleaner struct {
	transactor              gateway.Transactor
	todoQueriesGateway      gateway.TodoQueriesGateway
	todoCommandsGateway     gateway.TodoCommandsGateway
	todoListQueriesGateway  gateway.TodoListQueriesGateway
	todoListCommandsGateway gateway.TodoListCommandsGateway
}

func NewTodoSoftDeletedCleaner(
	transactor gateway.Transactor,
	todoQueriesGateway gateway.TodoQueriesGateway,
	todoCommandsGateway gateway.TodoCommandsGateway,
	todoListQueriesGateway gateway.TodoListQueriesGateway,
	todoListCommandsGateway gateway.TodoListCommandsGateway,
) *TodoSoftDeletedCleaner {
	return &TodoSoftDeletedCleaner{
		transactor:              transactor,
		todoQueriesGateway:      todoQueriesGateway,
		todoCommandsGateway:     todoCommandsGateway,
		todoListQueriesGateway:  todoListQueriesGateway,
		todoListCommandsGateway: todoListCommandsGateway,
	}
}

func (s *TodoSoftDeletedCleaner) Clean(
	ctx context.Context,
	in *input.TodoSoftDeletedCleaner,
) (*output.TodoSoftDeletedCleaner, error) {
	if in == nil {
		return nil, entity.NewInvalidParameter("todo soft deleted cleaner input is nil")
	}
	if in.AsOf.IsZero() {
		return nil, entity.NewInvalidParameter("as_of is required")
	}
	if in.RetentionDays <= 0 {
		return nil, entity.NewInvalidParameter("retention_days must be greater than zero")
	}
	if in.BatchSize <= 0 {
		return nil, entity.NewInvalidParameter("batch_size must be greater than zero")
	}
	if in.MaxBatches <= 0 {
		return nil, entity.NewInvalidParameter("max_batches must be greater than zero")
	}

	cutoff := in.AsOf.AddDate(0, 0, -in.RetentionDays)

	var (
		deletedTodoListCount int64
		deletedTodoCount     int64
		batchCount           int
		hasMore              bool
	)

	for batchCount < in.MaxBatches {
		var (
			selectedTodoListCount int
			deletedTodoLists      int64
			deletedChildTodos     int64
		)

		err := s.transactor.Transaction(ctx, func(txCtx context.Context) error {
			ids, err := s.todoListQueriesGateway.FindSoftDeletedTodoListIDs(ctx, cutoff, in.BatchSize)
			if err != nil {
				return fmt.Errorf("find soft deleted todo list ids: %w", err)
			}

			selectedTodoListCount = len(ids)
			if selectedTodoListCount == 0 {
				return nil
			}

			deletedChildTodos, err = s.todoCommandsGateway.HardDeleteTodosByTodoListIDs(txCtx, ids)
			if err != nil {
				return fmt.Errorf("hard delete todos by todo list ids: %w", err)
			}

			deletedTodoLists, err = s.todoListCommandsGateway.HardDeleteTodoListsByIDs(txCtx, ids)
			if err != nil {
				return fmt.Errorf("hard delete todo lists by ids: %w", err)
			}

			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("TodoSoftDeletedCleaner.Clean: clean todo list batch: %w", err)
		}

		if selectedTodoListCount == 0 {
			break
		}

		batchCount++
		deletedTodoListCount += deletedTodoLists
		deletedTodoCount += deletedChildTodos

		logutil.Info(ctx, "todo soft deleted cleaner todo list batch completed",
			zap.Int("batch_number", batchCount),
			zap.Int("selected_todo_lists", selectedTodoListCount),
			zap.Int64("deleted_todo_lists", deletedTodoLists),
			zap.Int64("deleted_child_todos", deletedChildTodos),
		)

		if selectedTodoListCount < in.BatchSize {
			break
		}

		if batchCount >= in.MaxBatches {
			hasMore = true
			break
		}

		if in.BatchSleep > 0 {
			select {
			case <-time.After(in.BatchSleep):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

	}

	for batchCount < in.MaxBatches {
		ids, err := s.todoQueriesGateway.FindSoftDeletedTodoIDs(ctx, cutoff, in.BatchSize)
		if err != nil {
			return nil, fmt.Errorf("TodoSoftDeletedCleaner.Clean: find soft deleted todo ids: %w", err)
		}

		if len(ids) == 0 {
			break
		}

		deletedTodos, err := s.todoCommandsGateway.HardDeleteTodosByIDs(ctx, ids)
		if err != nil {
			return nil, fmt.Errorf("TodoSoftDeletedCleaner.Clean: hard delete todos by ids: %w", err)
		}

		batchCount++
		deletedTodoCount += deletedTodos

		if len(ids) < in.BatchSize {
			break
		}

		if batchCount >= in.MaxBatches {
			hasMore = true
			break
		}

		if in.BatchSleep > 0 {
			select {
			case <-time.After(in.BatchSleep):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

	}

	return &output.TodoSoftDeletedCleaner{
		DeletedTodoListCount: deletedTodoListCount,
		DeletedTodoCount:     deletedTodoCount,
		BatchCount:           batchCount,
		HasMore:              hasMore,
	}, nil
}
