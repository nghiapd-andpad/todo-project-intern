// Package gateway defines the interfaces for interacting with external services or data sources.
package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
)

// filter options
type ListTodoListsOptions struct {
	NameSearch *string
	Offset     int
	Limit      int
}

type ListTodosOptions struct {
	Status      *entity.TodoStatus
	Priority    *entity.Priority
	TitleSearch *string
	Offset      int
	Limit       int
}

type TodoListPage struct {
	TodoLists []*entity.TodoList
	Total     int64
}

type TodoPage struct {
	Todos []*entity.Todo
	Total int64
}

type CreateTodoInput struct {
	Title       string
	Description *string
	Priority    *entity.Priority
	DueDate     *string // "2006-01-02"
	AssigneeID  *string
}

type UpdateTodoInput struct {
	Title       *string
	Description *string
	Status      *entity.TodoStatus
	Priority    *entity.Priority
	DueDate     *string
	AssigneeID  *string
}

type CreateTodoListInput struct {
	Parent      string
	DisplayName string
}

type UpdateTodoListInput struct {
	Name        string
	DisplayName *string
}

type TodoGateway interface {
	// TodoList operations
	GetTodoList(ctx context.Context, name string) (*entity.TodoList, error)
	ListTodoLists(ctx context.Context, parent string, opts ListTodoListsOptions) (*TodoListPage, error)
	CreateTodoList(ctx context.Context, input CreateTodoListInput) (*entity.TodoList, error)
	UpdateTodoList(ctx context.Context, input UpdateTodoListInput) (*entity.TodoList, error)
	DeleteTodoList(ctx context.Context, name string) error

	// Todo operations
	GetTodo(ctx context.Context, name string) (*entity.Todo, error)
	ListTodos(ctx context.Context, parent string, opts ListTodosOptions) (*TodoPage, error)
	CreateTodo(ctx context.Context, parent string, input CreateTodoInput) (*entity.Todo, error)
	UpdateTodo(ctx context.Context, name string, input UpdateTodoInput) (*entity.Todo, error)
	DeleteTodo(ctx context.Context, name string) error
}
