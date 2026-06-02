package input

import (
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

// TodoListFilter defines the scope of ListTodoLists.
type TodoListFilter string

const (
	// TodoListFilterAll returns lists owned by the user or containing todos assigned to the user.
	TodoListFilterAll TodoListFilter = "ALL"
	// TodoListFilterOwned returns only lists owned by the user.
	TodoListFilterOwned TodoListFilter = "OWNED"
	// TodoListFilterAssigned returns only lists containing todos assigned to the user.
	TodoListFilterAssigned TodoListFilter = "ASSIGNED"
)

type TodoListGetter struct {
	TodoListID  entity.TodoListID
	RequesterID entity.UserID
}

type TodoListLister struct {
	RequesterID entity.UserID
	// Filter defaults to TodoListFilterAll.
	Filter TodoListFilter
	// Optional filters
	NameSearch *string
	Offset     int
	Limit      int
}

type TodoListCreator struct {
	Name           string
	RequesterID    entity.UserID
	IdempotencyKey *string
}

type TodoListUpdater struct {
	TodoListID  entity.TodoListID
	RequesterID entity.UserID
	Name        *string
}

type TodoListDeleter struct {
	TodoListID  entity.TodoListID
	RequesterID entity.UserID
}
