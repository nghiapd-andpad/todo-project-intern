// Package output contains the output data structures for the todo usecases.
package output

import (
	"time"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
)

type TodoListOutput struct {
	ID          string
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TodoOutput struct {
	ID          string
	TodoListID  string
	Title       string
	Description *string
	Status      entity.TodoStatus
	Priority    entity.Priority
	DueDate     *time.Time
	CreatorID   string
	AssigneeID  *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TodoListPage struct {
	TodoLists []*TodoListOutput
	Total     int64
}

type TodoPage struct {
	Todos []*TodoOutput
	Total int64
}
