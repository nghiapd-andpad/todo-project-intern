package gateway

import (
	"context"

	"github.com/nghiaphunng18/todos/internal/domain/entity"
)

type TodoCommandsGateway interface {
	Create(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
	Update(ctx context.Context, todo *entity.Todo) (*entity.Todo, error)
	Delete(ctx context.Context, todoID entity.TodoID) error
}
