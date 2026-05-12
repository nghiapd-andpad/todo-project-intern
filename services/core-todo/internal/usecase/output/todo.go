// Package output defines the output data structures for use cases.
package output

import (
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/shared/pagination"
)

type TodoGetter struct {
	Todo *entity.Todo
}

type TodoLister struct {
	Page *pagination.Page[*entity.Todo]
}

type TodoCreator struct {
	Todo *entity.Todo
}

type TodoUpdater struct {
	Todo *entity.Todo
}

type TodoDeleter struct{}
