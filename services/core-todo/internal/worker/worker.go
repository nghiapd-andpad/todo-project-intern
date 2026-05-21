// Package worker wires and starts background jobs for the core-todo service.
package worker

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/worker/job"
)

type Worker struct {
	cfg       *config.Config
	scheduler gateway.Scheduler
	jobs      []job.CronJob
	logger    *zap.Logger
}

func NewWorker(
	cfg *config.Config,
	scheduler gateway.Scheduler,
	todoOverdueMarkerJob *job.TodoOverdueMarkerJob,
	zapLogger *zap.Logger,
) *Worker {
	if zapLogger == nil {
		zapLogger = zap.NewNop()
	}

	return &Worker{
		cfg:       cfg,
		scheduler: scheduler,
		jobs: []job.CronJob{
			todoOverdueMarkerJob,
		},
		logger: zapLogger.With(zap.String("component", "worker")),
	}
}

func (w *Worker) Start(ctx context.Context) error {
	if w.cfg == nil || !w.cfg.SchedulerEnabled {
		w.logger.Info("worker scheduler disabled")
		return nil
	}

	for _, j := range w.jobs {
		if err := w.scheduler.ScheduleCron(ctx, j.Name(), j.Cron(), j.Run); err != nil {
			return fmt.Errorf("register job %s: %w", j.Name(), err)
		}
	}

	w.scheduler.Start()

	w.logger.Info("worker scheduler started",
		zap.Int("job_count", len(w.jobs)),
	)

	return nil
}

func (w *Worker) Stop() error {
	if w.scheduler == nil {
		return nil
	}

	return w.scheduler.Stop()
}
