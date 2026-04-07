package entity

import "time"

type Todo struct {
	ID          TodoID
	TodoListID  *TodoListID
	Title       string
	Description *string
	Status      TodoStatus
	Priority    Priority
	DueDate     *time.Time
	AssigneeID  *UserID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTodo(title string) *Todo {
	now := time.Now()

	return &Todo{
		Title:     title,
		Status:    TodoStatusPending,
		Priority:  PriorityMedium,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (t *Todo) IsOverdue() bool {
	if t.DueDate == nil || t.Status == TodoStatusDone {
		return false
	}
	return time.Now().After(*t.DueDate)
}
