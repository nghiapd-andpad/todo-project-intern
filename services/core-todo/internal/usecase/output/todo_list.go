package output

import (
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/shared/pagination"
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
