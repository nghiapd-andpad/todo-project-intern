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
func CreateTodo(t *testing.T, db *gorm.DB, todoListID entity.TodoListID, title string) *entity.Todo {
	t.Helper()

	repo := persistence.NewTodoCommandsGateway(db)

	todo, err := repo.Create(context.Background(), &entity.Todo{
		TodoListID: todoListID,
		Title:      title,
		Status:     entity.TodoStatusPending,
		Priority:   entity.PriorityMedium,
	})
	if err != nil {
		t.Fatalf("fixture CreateTodo: %v", err)
	}

	return todo
}

func CreateTodoWithAssignee(t *testing.T, db *gorm.DB, todoListID entity.TodoListID, title string, assigneeID entity.UserID) *entity.Todo {
	t.Helper()

	repo := persistence.NewTodoCommandsGateway(db)

	todo, err := repo.Create(context.Background(), &entity.Todo{
		TodoListID: todoListID,
		Title:      title,
		Status:     entity.TodoStatusPending,
		Priority:   entity.PriorityMedium,
		AssigneeID: &assigneeID,
	})
	if err != nil {
		t.Fatalf("fixture CreateTodoWithAssignee: %v", err)
	}

	return todo
}
