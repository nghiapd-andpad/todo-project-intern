// Package cronjob provides gocron-based scheduler infrastructure.
package cronjob

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
)

const defaultLimitConcurrentJobs = 10

type Scheduler struct {
	cronScheduler gocron.Scheduler
	logger        *zap.Logger

	stopOnce sync.Once
	stopErr  error
}

var _ gateway.Scheduler = (*Scheduler)(nil)

func NewScheduler(zapLogger *zap.Logger) (*Scheduler, func(), error) {
	if zapLogger == nil {
		zapLogger = zap.NewNop()
	}

	l := newLogger(zapLogger)

	cronScheduler, err := gocron.NewScheduler(
		gocron.WithLogger(l),
		gocron.WithLimitConcurrentJobs(defaultLimitConcurrentJobs, gocron.LimitModeWait),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create gocron scheduler: %w", err)
	}

	s := &Scheduler{
		cronScheduler: cronScheduler,
		logger:        zapLogger.With(zap.String("component", "cronjob")),
	}

	cleanup := func() {
		if err := s.Stop(); err != nil {
			l.Error("failed to shutdown scheduler: %v", err)
		}
	}

	return s, cleanup, nil
}

func (s *Scheduler) ScheduleCron(
	ctx context.Context,
	name string,
	cronExpr string,
	task any,
	parameters ...any,
) error {
	if name == "" {
		return fmt.Errorf("job name is empty")
	}
	if cronExpr == "" {
		return fmt.Errorf("cron expression is empty for job %s", name)
	}

	_, err := s.cronScheduler.NewJob(
		gocron.CronJob(cronExpr, false),
		gocron.NewTask(task, parameters...),
		gocron.WithName(name),
		gocron.WithTags(name),
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
		return fmt.Errorf("register cronjob %s: %w", name, err)
	}

	s.logger.Info("cronjob registered",
		zap.String("job_name", name),
		zap.String("cron", cronExpr),
	)

	return nil
}

func (s *Scheduler) Start() {
	s.logger.Info("scheduler started")
	s.cronScheduler.Start()
}

func (s *Scheduler) Stop() error {
	s.stopOnce.Do(func() {
		if s.cronScheduler != nil {
			s.stopErr = s.cronScheduler.Shutdown()
		}
	})
	return s.stopErr
}
