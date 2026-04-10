//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/nghiaphunng18/todos/internal/handler/grpc/service"
	"github.com/nghiaphunng18/todos/internal/infra/persistence"
	"github.com/nghiaphunng18/todos/internal/usecase/todos"
	"gorm.io/gorm"
)

func InitializeTodoHandler(db *gorm.DB) *service.TodoHandler {
	wire.Build(
		persistence.WireSet,
		todos.WireSet,
		service.NewTodoHandler,
	)
	return &service.TodoHandler{}
}
