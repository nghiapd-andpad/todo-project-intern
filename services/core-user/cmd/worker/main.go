package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/di"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	logutil "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/utils/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app, cleanup, err := di.InitializeWorker(cfg)
	if err != nil {
		log.Fatalf("failed to initialize worker: %v", err)
	}
	defer cleanup()

	logutil.Info(ctx, "worker starting",
		zap.Bool("scheduler_enabled", cfg.SchedulerEnabled),
	)

	if err := app.Worker.Start(ctx); err != nil {
		logutil.Error(ctx, "failed to start worker", zap.Error(err))
		return
	}

	logutil.Info(ctx, "worker started and waiting for shutdown signal")

	<-ctx.Done()

	logutil.Info(context.Background(), "worker shutting down")

	if err := app.Worker.Stop(); err != nil {
		logutil.Error(context.Background(), "failed to stop worker", zap.Error(err))
	}

	logutil.Info(context.Background(), "worker stopped")
}
