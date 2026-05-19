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

type TodoOverdueMarker struct {
	todoQueriesGateway  gateway.TodoQueriesGateway
	todoCommandsGateway gateway.TodoCommandsGateway
}

func NewTodoOverdueMarker(
	todoQueriesGateway gateway.TodoQueriesGateway,
	todoCommandsGateway gateway.TodoCommandsGateway,
) *TodoOverdueMarker {
	return &TodoOverdueMarker{
		todoQueriesGateway:  todoQueriesGateway,
		todoCommandsGateway: todoCommandsGateway,
	}
}

func (s *TodoOverdueMarker) MarkOverdue(
	ctx context.Context,
	in *input.TodoOverdueMarker,
) (*output.TodoOverdueMarker, error) {
	if in == nil {
		return nil, entity.NewInvalidParameter("todo overdue marker input is nil")
	}
	if in.AsOf.IsZero() {
		return nil, entity.NewInvalidParameter("as_of is required")
	}
	if in.BatchSize <= 0 {
		return nil, entity.NewInvalidParameter("batch_size must be greater than zero")
	}
	if in.MaxBatches <= 0 {
		return nil, entity.NewInvalidParameter("max_batches must be greater than zero")
	}

	logutil.Info(ctx, "todo overdue marker started",
		zap.Time("as_of", in.AsOf),
		zap.Int("batch_size", in.BatchSize),
		zap.Int("max_batches", in.MaxBatches),
	)

	var markedCount int64
	var batchCount int
	var hasMore bool

	for batch := 0; batch < in.MaxBatches; batch++ {
		ids, err := s.todoQueriesGateway.FindOverdueTodoIDs(ctx, in.AsOf, in.BatchSize)
		if err != nil {
			return nil, fmt.Errorf("TodoOverdueMarker.MarkOverdue: find overdue todo ids: %w", err)
		}

		if len(ids) == 0 {
			hasMore = false
			break
		}

		affected, err := s.todoCommandsGateway.MarkOverdueByIDs(ctx, ids, in.AsOf)
		if err != nil {
			return nil, fmt.Errorf("TodoOverdueMarker.MarkOverdue: mark overdue todos by ids: %w", err)
		}

		batchCount++
		markedCount += affected

		logutil.Info(ctx, "todo overdue marker batch completed",
			zap.Int("batch_number", batchCount),
			zap.Int("selected_count", len(ids)),
			zap.Int64("marked_count", affected),
		)

		// If selected count is less than batch size, this was the last available batch.
		if len(ids) < in.BatchSize {
			hasMore = false
			break
		}

		// If this is the last allowed batch and we still selected a full batch,
		// there may be more overdue todos left for the next scheduler run.
		if batch == in.MaxBatches-1 {
			hasMore = true
		}

		if in.BatchSleep > 0 && batch < in.MaxBatches-1 {
			time.Sleep(in.BatchSleep)
		}
	}

	logutil.Info(ctx, "todo overdue marker completed",
		zap.Int64("marked_count", markedCount),
		zap.Int("batch_count", batchCount),
		zap.Bool("has_more", hasMore),
	)

	return &output.TodoOverdueMarker{
		MarkedCount: markedCount,
		BatchCount:  batchCount,
		HasMore:     hasMore,
	}, nil
}
