package main

import (
	"fmt"
	"log"
	"net"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/di"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
)

func main() {
	// Load config
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize App
	server, cleanup, err := di.InitializeApp(cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer cleanup()

	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("gRPC Server is running on port :%s\n", cfg.ServerPort)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
