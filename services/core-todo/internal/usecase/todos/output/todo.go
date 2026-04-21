package output

import "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"

type TodoGetter struct {
	Todo *entity.Todo
}

type TodoLister struct {
	Todos []*entity.Todo
	Total int64
}

type TodoCreator struct {
	Todo *entity.Todo
}

type TodoUpdater struct {
	Todo *entity.Todo
}

type TodoDeleter struct{}
