package entity

import "time"

type TodoID int64

type TodoStatus string

const (
	TodoStatusPending    TodoStatus = "PENDING"
	TodoStatusInProgress TodoStatus = "IN_PROGRESS"
	TodoStatusDone       TodoStatus = "DONE"
)

type Priority string

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
	PriorityUrgent Priority = "URGENT"
)

type Todo struct {
	ID          TodoID
	TodoListID  TodoListID
	Title       string
	Description *string
	Status      TodoStatus
	Priority    Priority
	AssigneeID  *UserID
	DueDate     *time.Time
	CreatorID   UserID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
