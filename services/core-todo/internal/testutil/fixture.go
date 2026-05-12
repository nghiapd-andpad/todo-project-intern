package testutil

import (
	"context"
	"testing"

	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence"
)

// CreateTodoList inserts a TodoList fixture and returns the created entity.
func CreateTodoList(t *testing.T, db *gorm.DB, name string, ownerID entity.UserID) *entity.TodoList {
	t.Helper()

	repo := persistence.NewTodoListCommandsGateway(db)
	todoList, err := repo.Create(context.Background(), &entity.TodoList{
		Name:    name,
		OwnerID: ownerID,
	})
	if err != nil {
		t.Fatalf("fixture CreateTodoList: %v", err)
	}
	return todoList
}

// CreateTodo inserts a Todo fixture and returns the created entity.
func CreateTodo(t *testing.T, db *gorm.DB, todoListID entity.TodoListID, title string, creatorID entity.UserID) *entity.Todo {
	t.Helper()

	repo := persistence.NewTodoCommandsGateway(db)
	todo, err := repo.Create(context.Background(), &entity.Todo{
		TodoListID: todoListID,
		Title:      title,
		Status:     entity.TodoStatusPending,
		Priority:   entity.PriorityMedium,
		// CreatorID:  creatorID,
	})
	if err != nil {
		t.Fatalf("fixture CreateTodo: %v", err)
	}
	return todo
}
