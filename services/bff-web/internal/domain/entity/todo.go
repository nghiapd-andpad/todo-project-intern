package entity

import "time"

type TodoList struct {
	Name        string // users/{user_id}/todo-lists/{list_id}
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Todo struct {
	Name        string // users/{user_id}/todo-lists/{list_id}/todos/{todo_id}
	Title       string
	Description *string
	Status      TodoStatus
	Priority    Priority
	DueDate     *time.Time
	CreatorID   string
	AssigneeID  *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

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
