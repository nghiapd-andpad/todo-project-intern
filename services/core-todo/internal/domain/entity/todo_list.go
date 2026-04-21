package entity

import "time"

type TodoListID int64

type TodoList struct {
	ID        TodoListID
	Name      string
	OwnerID   UserID
	CreatedAt time.Time
	UpdatedAt time.Time
}
