// Package input defines the input data structures for gateways.
package input

import (
	"time"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

type ListTodosOptions struct {
	TodoListID   entity.TodoListID
	AssigneeOnly *entity.UserID
	Status       *entity.TodoStatus
	Priority     *entity.Priority
	TitleSearch  *string
	Offset       int
	Limit        int

	// cursor pagination
	CursorCreatedAt *time.Time
	CursorID        *entity.TodoID
}
