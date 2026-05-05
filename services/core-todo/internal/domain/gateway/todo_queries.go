// Package gateway defines the interfaces for data access and external service communication.
//
//go:generate mockgen -destination=mock/todo_queries_mock.go -source=todo_queries.go -package mock
package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

type ListTodosOptions struct {
	TodoListID  *entity.TodoListID
	Status      *entity.TodoStatus
	Priority    *entity.Priority
	CreatorID   *entity.UserID
	AssigneeID  *entity.UserID
	TitleSearch *string
	Offset      int
	Limit       int
}

func (o *ListTodosOptions) Validate() error {
	if o.Limit < 0 {
		return entity.NewInvalidParameter("limit must be non-negative")
	}
	if o.Limit > 100 {
		return entity.NewInvalidParameter("limit must not exceed 100")
	}
	if o.Offset < 0 {
		return entity.NewInvalidParameter("offset must be non-negative")
	}
	if o.TitleSearch != nil && len(*o.TitleSearch) > 255 {
		return entity.NewInvalidParameter("title_search too long")
	}
	return nil
}

func (o *ListTodosOptions) SetDefaults() {
	if o.Limit == 0 {
		o.Limit = 20
	}
}

type TodoQueriesGateway interface {
	Get(ctx context.Context, todoID entity.TodoID) (*entity.Todo, error)
	List(ctx context.Context, opts ListTodosOptions) ([]*entity.Todo, int64, error)
}
