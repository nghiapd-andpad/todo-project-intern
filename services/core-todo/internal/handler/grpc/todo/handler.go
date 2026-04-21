package todo

import (
	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

func NewGRPCServer(cfg *config.Config, handler *TodoHandler) *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor()),
	)
	todov1.RegisterTodosServiceServer(s, handler)
	reflection.Register(s)
	return s
}
