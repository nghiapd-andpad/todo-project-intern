package job

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	logutil "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/utils/logger"
)

const todoOverdueMarkerJobName = "todo_overdue_marker"

type TodoOverdueMarkerJob struct {
	marker usecase.TodoOverdueMarker
	locker gateway.DistributedLocker
	cfg    *config.Config
	logger *zap.Logger
}

func NewTodoOverdueMarkerJob(
	marker usecase.TodoOverdueMarker,
	locker gateway.DistributedLocker,
	cfg *config.Config,
	zapLogger *zap.Logger,
) *TodoOverdueMarkerJob {
	if zapLogger == nil {
		zapLogger = zap.NewNop()
	}

	return &TodoOverdueMarkerJob{
		marker: marker,
		locker: locker,
		cfg:    cfg,
		logger: zapLogger.With(
			zap.String("component", "cronjob"),
			zap.String("job_name", todoOverdueMarkerJobName),
		),
	}
}

func (j *TodoOverdueMarkerJob) Name() string {
	return todoOverdueMarkerJobName
}

func (j *TodoOverdueMarkerJob) Cron() string {
	if j.cfg == nil {
		return ""
	}
	return j.cfg.TodoOverdueMarkerCron
}

func (j *TodoOverdueMarkerJob) Run() error {
	start := time.Now()

	ctx := logutil.ToContext(context.Background(), j.logger)

	logutil.Info(ctx, "todo overdue marker job started")

	acquired, release, err := j.locker.TryLock(
		ctx,
		j.cfg.TodoOverdueMarkerLockKey,
		j.cfg.TodoOverdueMarkerLockTTL,
	)
	if err != nil {
		return fmt.Errorf("acquire todo overdue marker lock: %w", err)
	}

	if !acquired {
		logutil.Info(ctx, "todo overdue marker job skipped because lock is held")
		return nil
	}
	defer release()

	out, err := j.marker.MarkOverdue(ctx, &input.TodoOverdueMarker{
		AsOf:       time.Now().UTC(),
		BatchSize:  j.cfg.TodoOverdueMarkerBatchSize,
		MaxBatches: j.cfg.TodoOverdueMarkerMaxBatches,
		BatchSleep: j.cfg.TodoOverdueMarkerBatchSleep,
	})
	if err != nil {
		return fmt.Errorf("mark overdue todos: %w", err)
	}

	logutil.Info(ctx, "todo overdue marker job completed",
		zap.Int64("marked_count", out.MarkedCount),
		zap.Int("batch_count", out.BatchCount),
		zap.Bool("has_more", out.HasMore),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}
