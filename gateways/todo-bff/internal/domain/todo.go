package domain

import "time"

type TodoStatus string
type Priority string

const (
	StatusPending    TodoStatus = "PENDING"
	StatusInProgress TodoStatus = "IN_PROGRESS"
	StatusDone       TodoStatus = "DONE"
)

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
)

type Todo struct {
	ID          string
	Title       string
	Description *string
	Status      TodoStatus
	Priority    Priority
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
