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
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/output"
)

type TodoCreatorUsecase interface {
	Create(ctx context.Context, in *input.TodoCreator) (*output.TodoCreator, error)
}

type TodoGetterUsecase interface {
	Get(ctx context.Context, in *input.TodoGetter) (*output.TodoGetter, error)
}

type TodoListerUsecase interface {
	List(ctx context.Context, in *input.TodoLister) (*output.TodoLister, error)
}

type TodoUpdaterUsecase interface {
	Update(ctx context.Context, in *input.TodoUpdater) (*output.TodoUpdater, error)
}

type TodoDeleterUsecase interface {
	Delete(ctx context.Context, in *input.TodoDeleter) (*output.TodoDeleter, error)
}

type TodoListCreatorUsecase interface {
	Create(ctx context.Context, in *input.TodoListCreator) (*output.TodoListCreator, error)
}

type TodoListGetterUsecase interface {
	Get(ctx context.Context, in *input.TodoListGetter) (*output.TodoListGetter, error)
}

type TodoListListerUsecase interface {
	List(ctx context.Context, in *input.TodoListLister) (*output.TodoListLister, error)
}

type TodoListUpdaterUsecase interface {
	Update(ctx context.Context, in *input.TodoListUpdater) (*output.TodoListUpdater, error)
}

type TodoListDeleterUsecase interface {
	Delete(ctx context.Context, in *input.TodoListDeleter) (*output.TodoListDeleter, error)
}

type TodoHandler struct {
	todov1.UnimplementedTodosServiceServer

	todoCreator TodoCreatorUsecase
	todoGetter  TodoGetterUsecase
	todoLister  TodoListerUsecase
	todoUpdater TodoUpdaterUsecase
	todoDeleter TodoDeleterUsecase

	todoListCreator TodoListCreatorUsecase
	todoListGetter  TodoListGetterUsecase
	todoListLister  TodoListListerUsecase
	todoListUpdater TodoListUpdaterUsecase
	todoListDeleter TodoListDeleterUsecase
}

func NewTodoHandler(
	todoCreator TodoCreatorUsecase,
	todoGetter TodoGetterUsecase,
	todoLister TodoListerUsecase,
	todoUpdater TodoUpdaterUsecase,
	todoDeleter TodoDeleterUsecase,

	todoListCreator TodoListCreatorUsecase,
	todoListGetter TodoListGetterUsecase,
	todoListLister TodoListListerUsecase,
	todoListUpdater TodoListUpdaterUsecase,
	todoListDeleter TodoListDeleterUsecase,
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
