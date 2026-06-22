// Package worker wires and starts background consumers for the core-user service.
package worker

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/infra/rabbitmq"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/worker/job"
)

type Worker struct {
	cfg *config.Config

	scheduler gateway.Scheduler
	jobs      []job.CronJob

	notificationConsumer *rabbitmq.NotificationConsumer
	emailConsumer        *rabbitmq.EmailConsumer

	logger *zap.Logger
}

func NewWorker(
	cfg *config.Config,
	scheduler gateway.Scheduler,

	notificationConsumer *rabbitmq.NotificationConsumer,
	emailConsumer *rabbitmq.EmailConsumer,
	outboxPublisherJob *job.OutboxPublisherJob,

	zapLogger *zap.Logger,
) *Worker {
	if zapLogger == nil {
		zapLogger = zap.NewNop()
	}

	return &Worker{
		cfg:                  cfg,
		scheduler:            scheduler,
		notificationConsumer: notificationConsumer,
		emailConsumer:        emailConsumer,
		jobs: []job.CronJob{
			outboxPublisherJob,
		},
		logger: zapLogger.With(zap.String("component", "worker")),
	}
}

func (w *Worker) Start(ctx context.Context) error {
	if err := w.notificationConsumer.Start(ctx); err != nil {
		return fmt.Errorf("start notification consumer: %w", err)
	}

	if err := w.emailConsumer.Start(ctx); err != nil {
		return fmt.Errorf("start email consumer: %w", err)
	}

	if w.cfg != nil &&
		w.cfg.SchedulerEnabled {

		for _, j := range w.jobs {
			if err := w.scheduler.ScheduleCron(
				ctx,
				j.Name(),
				j.Cron(),
				j.Run,
			); err != nil {
				return err
			}
		}

		w.scheduler.Start()
	}

	w.logger.Info("worker started")

	return nil
}

func (w *Worker) Stop() error {
	if w.scheduler != nil {
		return w.scheduler.Stop()
	}

	return nil
}
