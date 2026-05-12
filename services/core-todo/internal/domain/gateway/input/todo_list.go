package input

import "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"

type ListTodoListsOptions struct {
	OwnerID    *entity.UserID
	AssigneeID *entity.UserID
	NameSearch *string
	Offset     int
	Limit      int
}
