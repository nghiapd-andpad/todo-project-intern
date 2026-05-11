package output

import "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"

type TodoListGetter struct {
	TodoList *entity.TodoList
}

type TodoListLister struct {
	TodoLists []*entity.TodoList
	Total     int64
}

type TodoListCreator struct {
	TodoList *entity.TodoList
}

type TodoListUpdater struct {
	TodoList *entity.TodoList
}

type TodoListDeleter struct{}
