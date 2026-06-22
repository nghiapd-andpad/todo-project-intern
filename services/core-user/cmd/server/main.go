package main

import (
	"context"
	"log"
	"net"
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

	app, cleanup, err := di.InitializeServer(cfg)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}
	defer cleanup()

	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		app.Logger.Error("failed to listen",
			zap.String("port", cfg.ServerPort),
			zap.Error(err),
		)
		return
	}

	app.Logger.Info("grpc server starting",
		zap.String("port", cfg.ServerPort),
		zap.String("app_env", cfg.AppEnv),
	)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- app.GRPCServer.Serve(lis)
	}()

	select {
	case err := <-serverErr:
		logutil.Error(context.Background(), "grpc server exited unexpectedly", zap.Error(err))
	case <-ctx.Done():
		logutil.Info(context.Background(), "shutdown signal received — stopping grpc server")
		// GracefulStop waits for in-flight RPCs to complete.
		app.GRPCServer.GracefulStop()
		logutil.Info(context.Background(), "grpc server stopped")
	}
}
