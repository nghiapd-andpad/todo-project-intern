package main

import (
	"fmt"
	"log"
	"net"

	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/user-service/di"
	"github.com/nghiapd-andpad/todo-project-intern/services/user-service/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/user-service/internal/infra/persistence"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize Database & Auto-migrate
	db, err := persistence.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Use Google Wire to initialize the UserHandler
	userHandler, err := di.InitializeUserHandler(db, cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// Setup gRPC Server
	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	userv1.RegisterUserServiceServer(s, userHandler)
	reflection.Register(s)

	fmt.Printf("gRPC Server is running on port :%s...\n", cfg.ServerPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
