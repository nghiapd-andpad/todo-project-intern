package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

type TodoQueriesGateway interface {
	Get(ctx context.Context, todoID entity.TodoID) (*entity.Todo, error)
	List(ctx context.Context) ([]*entity.Todo, error)
}
