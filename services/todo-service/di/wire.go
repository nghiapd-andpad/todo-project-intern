//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/handler/grpc/todo"
	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/usecase/todos"
	"gorm.io/gorm"
)

func InitializeTodoHandler(db *gorm.DB) *todo.TodoHandler {
	wire.Build(
		persistence.WireSet,
		todos.WireSet,
		todo.NewTodoHandler,
	)
	return &todo.TodoHandler{}
}
