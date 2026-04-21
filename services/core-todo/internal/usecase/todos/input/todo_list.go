package input

import (
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
)

type TodoListGetter struct {
	ID entity.TodoListID
}

type TodoListLister struct {
	Opts gateway.ListTodoListsOptions
}

type TodoListCreator struct {
	Name    string
	OwnerID entity.UserID
}

type TodoListUpdater struct {
	ID   entity.TodoListID
	Name *string
}

type TodoListDeleter struct {
	ID entity.TodoListID
}
