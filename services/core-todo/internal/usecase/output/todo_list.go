package output

import (
	"github.com/nghiapd-andpad/todo-project-intern/pkg/pagination"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

type TodoListGetter struct {
	TodoList *entity.TodoList
}

type TodoListLister struct {
	Page *pagination.Page[*entity.TodoList]
}

type TodoListCreator struct {
	TodoList *entity.TodoList
}

type TodoListUpdater struct {
	TodoList *entity.TodoList
}

type TodoListDeleter struct{}

type TodoOverdueMarker struct {
	MarkedCount int64
	BatchCount  int
	HasMore     bool
}
