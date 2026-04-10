package main

import (
	"fmt"
	"log"
	"net"

	"github.com/nghiaphunng18/todos/di" // Import folder di nơi chứa wire_gen.go
	todov1 "github.com/nghiaphunng18/todos/gen/todo/v1"
	"github.com/nghiaphunng18/todos/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Use Google Wire to initialize the TodoHandler
	todoHandler := di.InitializeTodoHandler(db)

	// Setup gRPC Server
	lis, err := net.Listen("tcp", ":"+cfg.ServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	todov1.RegisterTodosServiceServer(s, todoHandler)
	reflection.Register(s)

	fmt.Printf("gRPC Server is running on port :%s...\n", cfg.ServerPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
