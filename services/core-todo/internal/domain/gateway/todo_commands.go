package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

type TodoCommandsGateway interface {
	Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
	Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
	Delete(ctx context.Context, todoID entity.TodoID) error
}
