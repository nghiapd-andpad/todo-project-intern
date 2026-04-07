package main

import (
	"fmt"
	"log"
	"net"

	todov1 "github.com/nghiaphunng18/todos/gen/todo/v1"
	handler "github.com/nghiaphunng18/todos/internal/handler/grpc/service"
	"github.com/nghiaphunng18/todos/internal/usecase/todos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC Server
	s := grpc.NewServer()

	todoCreator := todos.NewTodoCreator(nil)
	todoGetter := todos.NewTodoGetter(nil)
	todoListReader := todos.NewTodoListReader(nil)
	todoUpdater := todos.NewTodoUpdater(nil, nil)
	todoDeleter := todos.NewTodoDeleter(nil)

	// Create TodoHandler with use cases
	todoHandler := handler.NewTodoHandler(
		todoCreator,
		todoGetter,
		todoListReader,
		todoUpdater,
		todoDeleter,
	)

	// Register TodoHandler to gRPC Server
	todov1.RegisterTodosServiceServer(s, todoHandler)

	// Register reflection service on gRPC server
	reflection.Register(s)

	fmt.Println("gRPC Server is running on port :50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
