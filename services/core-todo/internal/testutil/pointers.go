package testutil

import "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"

func StrPtr(s string) *string {
	return &s
}

func UserIDPtr(id entity.UserID) *entity.UserID {
	return &id
}

func TodoListIDPtr(id entity.TodoListID) *entity.TodoListID {
	return &id
}

func TodoStatusPtr(s entity.TodoStatus) *entity.TodoStatus {
	return &s
}
