// Package gateway defines the interfaces for data access and external service communication.
//
//go:generate mockgen -destination=mock/todo_queries_mock.go -source=todo_queries.go -package mock
package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
)

type TodoQueriesGateway interface {
	Get(ctx context.Context, todoID entity.TodoID, todoListID entity.TodoListID) (*entity.Todo, error)
	List(ctx context.Context, opts *input.ListTodosOptions) ([]*entity.Todo, int64, error)
}
