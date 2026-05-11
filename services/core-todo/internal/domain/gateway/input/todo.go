// Package input defines the input data structures for gateways.
package input

import "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"

type ListTodosOptions struct {
	TodoListID  *entity.TodoListID
	Status      *entity.TodoStatus
	Priority    *entity.Priority
	CreatorID   *entity.UserID
	AssigneeID  *entity.UserID
	TitleSearch *string
	Offset      int
	Limit       int
}
