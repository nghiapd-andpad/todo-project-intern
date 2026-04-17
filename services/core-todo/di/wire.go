//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	todo "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/todo"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos"
)

func InitializeTodoHandler(db *gorm.DB) *todo.TodoHandler {
	wire.Build(
		persistence.WireSet,
		todos.WireSet,
		todo.NewTodoHandler,
	)
	return &todo.TodoHandler{}
}
