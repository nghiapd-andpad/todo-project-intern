package cronjob

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
	logutil "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/utils/logger"
)

const todoOverdueMarkerJobName = "todo_overdue_marker"

type Scheduler struct {
	cronScheduler gocron.Scheduler
	marker        usecase.TodoOverdueMarker
	locker        gateway.DistributedLocker
	cfg           *config.Config
	logger        *zap.Logger
}

func NewScheduler(
	cronScheduler gocron.Scheduler,
	marker usecase.TodoOverdueMarker,
	locker gateway.DistributedLocker,
	cfg *config.Config,
	zapLogger *zap.Logger,
) (*Scheduler, error) {
	if zapLogger == nil {
		zapLogger = zap.NewNop()
	}

	s := &Scheduler{
		cronScheduler: cronScheduler,
		marker:        marker,
		locker:        locker,
		cfg:           cfg,
		logger: zapLogger.With(
			zap.String("component", "cronjob"),
		),
	}

	if cfg != nil && cfg.SchedulerEnabled {
		if err := s.registerTodoOverdueMarkerJob(); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Scheduler) Start() {
	if s.cfg == nil || !s.cfg.SchedulerEnabled {
		s.logger.Info("scheduler disabled")
		return
	}

	s.logger.Info("scheduler started")
	s.cronScheduler.Start()
}

func (s *Scheduler) Stop() error {
	if s.cronScheduler == nil {
		return nil
	}

	return s.cronScheduler.Shutdown()
}

func (s *Scheduler) registerTodoOverdueMarkerJob() error {
	if s.cfg.TodoOverdueMarkerCron == "" {
		return fmt.Errorf("todo overdue marker cron is empty")
	}

	_, err := s.cronScheduler.NewJob(
		gocron.CronJob(s.cfg.TodoOverdueMarkerCron, false),
		gocron.NewTask(s.runTodoOverdueMarkerJob),
		gocron.WithName(todoOverdueMarkerJobName),
		gocron.WithTags(todoOverdueMarkerJobName),
		gocron.WithEventListeners(
			gocron.AfterJobRuns(func(_ uuid.UUID, jobName string) {
				s.logger.Info("cronjob finished",
					zap.String("job_name", jobName),
				)
			}),
			gocron.AfterJobRunsWithError(func(_ uuid.UUID, jobName string, err error) {
				s.logger.Error("cronjob failed",
					zap.String("job_name", jobName),
					zap.Error(err),
				)
			}),
			gocron.AfterJobRunsWithPanic(func(_ uuid.UUID, jobName string, r any) {
				s.logger.Error("cronjob panicked",
					zap.String("job_name", jobName),
					zap.Any("panic", r),
				)
			}),
		),
	)
	if err != nil {
		return fmt.Errorf("register todo overdue marker job: %w", err)
	}

	s.logger.Info("cronjob registered",
		zap.String("job_name", todoOverdueMarkerJobName),
		zap.String("cron", s.cfg.TodoOverdueMarkerCron),
	)

	return nil
}

func (s *Scheduler) runTodoOverdueMarkerJob() error {
	start := time.Now()

	ctx := logutil.ToContext(context.Background(), s.logger.With(
		zap.String("job_name", todoOverdueMarkerJobName),
	))

	logutil.Info(ctx, "todo overdue marker job started")

	acquired, release, err := s.locker.TryLock(
		ctx,
		s.cfg.TodoOverdueMarkerLockKey,
		s.cfg.TodoOverdueMarkerLockTTL,
	)
	if err != nil {
		return fmt.Errorf("acquire todo overdue marker lock: %w", err)
	}

	if !acquired {
		logutil.Info(ctx, "todo overdue marker job skipped because lock is held")
		return nil
	}
	defer release()

	out, err := s.marker.MarkOverdue(ctx, &input.TodoOverdueMarker{
		AsOf:       time.Now().UTC(),
		BatchSize:  s.cfg.TodoOverdueMarkerBatchSize,
		MaxBatches: s.cfg.TodoOverdueMarkerMaxBatches,
		BatchSleep: s.cfg.TodoOverdueMarkerBatchSleep,
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
