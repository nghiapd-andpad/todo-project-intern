// Package input defines the input data structures for use cases.
package input

import (
	"time"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

type TodoGetter struct {
	TodoID      entity.TodoID
	TodoListID  entity.TodoListID
	RequesterID entity.UserID
}

type TodoLister struct {
	TodoListID  entity.TodoListID
	RequesterID entity.UserID
	// Optional filters
	Status      *entity.TodoStatus
	Priority    *entity.Priority
	TitleSearch *string
	Offset      int
	Limit       int
}

type TodoCreator struct {
	TodoListID     entity.TodoListID
	RequesterID    entity.UserID
	Title          string
	Description    *string
	Priority       entity.Priority
	DueDate        *time.Time
	AssigneeID     *entity.UserID
	IdempotencyKey *string
}

type UpdateTodoFields struct {
	Title       *string
	Description *string
	Status      *entity.TodoStatus
	Priority    *entity.Priority
	DueDate     *time.Time
	AssigneeID  *entity.UserID
}

type TodoUpdater struct {
	TodoID      entity.TodoID
	TodoListID  entity.TodoListID
	RequesterID entity.UserID
	Fields      UpdateTodoFields
}

type TodoDeleter struct {
	TodoID      entity.TodoID
	TodoListID  entity.TodoListID
	RequesterID entity.UserID
}

type TodoOverdueMarker struct {
	AsOf       time.Time
	BatchSize  int
	MaxBatches int
	BatchSleep time.Duration
}

type TodoSoftDeletedCleaner struct {
	AsOf          time.Time
	RetentionDays int
	BatchSize     int
	MaxBatches    int
	BatchSleep    time.Duration
}
