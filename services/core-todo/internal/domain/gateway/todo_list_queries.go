//go:generate mockgen -destination=mock/todo_list_queries_mock.go -source=todo_list_queries.go -package mock

package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

type ListTodoListsOptions struct {
	OwnerID    *entity.UserID
	NameSearch *string
	Offset     int
	Limit      int
}

type TodoListQueriesGateway interface {
	Get(ctx context.Context, todoListID entity.TodoListID) (*entity.TodoList, error)
	List(ctx context.Context, opts ListTodoListsOptions) ([]*entity.TodoList, int64, error)
}
