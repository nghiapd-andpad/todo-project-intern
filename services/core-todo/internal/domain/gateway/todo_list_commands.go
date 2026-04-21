package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

type TodoListCommandsGateway interface {
	Create(ctx context.Context, todoList *entity.TodoList) (*entity.TodoList, error)
	Update(ctx context.Context, todoList *entity.TodoList) (*entity.TodoList, error)
	Delete(ctx context.Context, todoListID entity.TodoListID) error
}
