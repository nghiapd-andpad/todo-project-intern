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
