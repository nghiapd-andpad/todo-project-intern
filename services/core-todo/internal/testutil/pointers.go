package testutil

import (
	"time"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

func StrPtr(s string) *string {
	return &s
}

func UserIDPtr(id entity.UserID) *entity.UserID {
	return &id
}

func TimePtr(v time.Time) *time.Time {
	return &v
}

func TodoListIDPtr(id entity.TodoListID) *entity.TodoListID {
	return &id
}

func TodoStatusPtr(s entity.TodoStatus) *entity.TodoStatus {
	return &s
}

func PriorityPtr(priority entity.Priority) *entity.Priority {
	return &priority
}
