package main

import (
	"context"
	"log"
	"net"

	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/di"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	logutil "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/utils/logger"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app, cleanup, err := di.InitializeApp(cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer cleanup()

	app.Scheduler.Start()

	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		logutil.Error(ctx, "failed to listen", zap.Error(err))
		return
	}

	logutil.Info(ctx, "grpc server started",
		zap.String("port", cfg.ServerPort),
		zap.String("app_env", cfg.AppEnv),
		zap.String("log_level", cfg.LogLevel),
		zap.String("log_format", cfg.LogFormat),
		zap.Bool("scheduler_enabled", cfg.SchedulerEnabled),
	)

	if err := app.GRPCServer.Serve(lis); err != nil {
		logutil.Error(ctx, "failed to serve grpc server", zap.Error(err))
		return
	}
}
