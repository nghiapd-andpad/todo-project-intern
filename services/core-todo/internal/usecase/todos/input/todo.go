package input

import (
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
)

type TodoGetter struct {
	ID entity.TodoID
}

type TodoLister struct {
	Opts gateway.ListTodosOptions
}

type TodoCreator struct {
	TodoListID  entity.TodoListID
	Title       string
	Description *string
	Priority    entity.Priority
	DueDate     *string // "2006-01-02" — parse ở usecase
	AssigneeID  *entity.UserID
	CreatorID   entity.UserID
}

type TodoUpdater struct {
	ID          entity.TodoID
	Title       *string
	Description *string
	Status      *entity.TodoStatus
	Priority    *entity.Priority
	DueDate     *string
	AssigneeID  *entity.UserID
}

type TodoDeleter struct {
	ID entity.TodoID
}
