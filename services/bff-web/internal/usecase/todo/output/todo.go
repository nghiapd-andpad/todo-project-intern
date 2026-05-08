// Package output contains the output data structures for the todo usecases.
package output

import "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"

type TodoListPage struct {
	TodoLists []*entity.TodoList
	Total     int64
}

type TodoPage struct {
	Todos []*entity.Todo
	Total int64
}
