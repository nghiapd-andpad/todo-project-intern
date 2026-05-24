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

const todoSoftDeletedCleanupJobName = "todo_soft_deleted_cleanup"

type TodoSoftDeletedCleanupJob struct {
	cleaner usecase.TodoSoftDeletedCleaner
	locker  gateway.DistributedLocker
	cfg     *config.Config
	logger  *zap.Logger
}

func NewTodoSoftDeletedCleanupJob(
	cleaner usecase.TodoSoftDeletedCleaner,
	locker gateway.DistributedLocker,
	cfg *config.Config,
	zapLogger *zap.Logger,
) *TodoSoftDeletedCleanupJob {
	if zapLogger == nil {
		zapLogger = zap.NewNop()
	}

	return &TodoSoftDeletedCleanupJob{
		cleaner: cleaner,
		locker:  locker,
		cfg:     cfg,
		logger: zapLogger.With(
			zap.String("component", "cronjob"),
			zap.String("job_name", todoSoftDeletedCleanupJobName),
		),
	}
}

func (j *TodoSoftDeletedCleanupJob) Name() string {
	return todoSoftDeletedCleanupJobName
}

func (j *TodoSoftDeletedCleanupJob) Cron() string {
	if j.cfg == nil {
		return ""
	}
	return j.cfg.TodoSoftDeletedCleanupCron
}

func (j *TodoSoftDeletedCleanupJob) Run() error {
	start := time.Now()

	ctx := logutil.ToContext(context.Background(), j.logger)

	logutil.Info(ctx, "todo soft deleted cleanup job started")

	acquired, release, err := j.locker.TryLock(
		ctx,
		j.cfg.TodoSoftDeletedCleanupLockKey,
		j.cfg.TodoSoftDeletedCleanupLockTTL,
	)
	if err != nil {
		return fmt.Errorf("acquire todo soft deleted cleanup lock: %w", err)
	}

	if !acquired {
		logutil.Info(ctx, "todo soft deleted cleanup job skipped because lock is held")
		return nil
	}
	defer release()

	out, err := j.cleaner.Clean(ctx, &input.TodoSoftDeletedCleaner{
		AsOf:          time.Now().UTC(),
		RetentionDays: j.cfg.TodoSoftDeletedCleanupRetentionDays,
		BatchSize:     j.cfg.TodoSoftDeletedCleanupBatchSize,
		MaxBatches:    j.cfg.TodoSoftDeletedCleanupMaxBatches,
		BatchSleep:    j.cfg.TodoSoftDeletedCleanupBatchSleep,
	})
	if err != nil {
		return fmt.Errorf("clean soft deleted todos: %w", err)
	}

	logutil.Info(ctx, "todo soft deleted cleanup job completed",
		zap.Int64("deleted_todo_lists", out.DeletedTodoListCount),
		zap.Int64("deleted_todos", out.DeletedTodoCount),
		zap.Int("batch_count", out.BatchCount),
		zap.Bool("has_more", out.HasMore),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}
