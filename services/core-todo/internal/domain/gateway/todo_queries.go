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

type TodoQueriesGateway interface {
	Get(ctx context.Context, todoID entity.TodoID) (*entity.Todo, error)
	List(ctx context.Context, opts ListTodosOptions) ([]*entity.Todo, int64, error)
}
