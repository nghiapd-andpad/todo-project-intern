package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/gateways/todo-bff/internal/domain"
)

type TodoGateway interface {
	CreateTodo(ctx context.Context, title string, desc *string) (*domain.Todo, error)
	GetTodo(ctx context.Context, name string) (*domain.Todo, error)
	ListTodos(ctx context.Context, parent string) ([]*domain.Todo, error)
}
