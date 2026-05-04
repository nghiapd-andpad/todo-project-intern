// Package todo contains gRPC handlers for todo service.
package todo

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos"
)

type TodoHandler struct {
	todov1.UnimplementedTodosServiceServer

	// Todo usecases
	todoCreator todos.TodoCreator
	todoGetter  todos.TodoGetter
	todoLister  todos.TodoLister
	todoUpdater todos.TodoUpdater
	todoDeleter todos.TodoDeleter

	// TodoList usecases
	todoListCreator todos.TodoListCreator
	todoListGetter  todos.TodoListGetter
	todoListLister  todos.TodoListLister
	todoListUpdater todos.TodoListUpdater
	todoListDeleter todos.TodoListDeleter
}

func NewTodoHandler(
	todoCreator todos.TodoCreator,
	todoGetter todos.TodoGetter,
	todoLister todos.TodoLister,
	todoUpdater todos.TodoUpdater,
	todoDeleter todos.TodoDeleter,
	todoListCreator todos.TodoListCreator,
	todoListGetter todos.TodoListGetter,
	todoListLister todos.TodoListLister,
	todoListUpdater todos.TodoListUpdater,
	todoListDeleter todos.TodoListDeleter,
) *TodoHandler {
	return &TodoHandler{
		todoCreator:     todoCreator,
		todoGetter:      todoGetter,
		todoLister:      todoLister,
		todoUpdater:     todoUpdater,
		todoDeleter:     todoDeleter,
		todoListCreator: todoListCreator,
		todoListGetter:  todoListGetter,
		todoListLister:  todoListLister,
		todoListUpdater: todoListUpdater,
		todoListDeleter: todoListDeleter,
	}
}

func NewGRPCServer(cfg *config.Config, handler *TodoHandler) (*grpc.Server, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("create validator: %w", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			auth.UnaryServerInterceptor(),
			func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
				if msg, ok := req.(proto.Message); ok {
					if err := validator.Validate(msg); err != nil {
						return nil, status.Error(codes.InvalidArgument, err.Error())
					}
				}
				return handler(ctx, req)
			},
		),
	)

	todov1.RegisterTodosServiceServer(s, handler)
	reflection.Register(s)
	return s, nil
}

func ProvideGRPCServer(cfg *config.Config, handler *TodoHandler) (*grpc.Server, func(), error) {
	s, err := NewGRPCServer(cfg, handler)
	if err != nil {
		return nil, nil, err
	}

	return s, func() {}, nil
}
