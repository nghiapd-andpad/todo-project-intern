package mapper

import (
	"fmt"

	todov1 "github.com/nghiaphunng18/todos/gen/todo/v1"
	"github.com/nghiaphunng18/todos/internal/domain/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TodoToPb(ent *entity.Todo) *todov1.Todo {
	if ent == nil {
		return nil
	}

	// Create Resource Name: users/{user_id}/todo-lists/{list_id}/todos/{todo_id}
	name := fmt.Sprintf("users/%d/todo-lists/%d/todos/%d", 1, 1, ent.ID)

	return &todov1.Todo{
		Name:        name,
		Title:       ent.Title,
		Description: *ent.Description,
		Status:      TodoStatusToPb(ent.Status),
		CreatedAt:   timestamppb.New(ent.CreatedAt),
		UpdatedAt:   timestamppb.New(ent.UpdatedAt),
	}
}

func TodoStatusToPb(status entity.TodoStatus) todov1.TodoStatus {
	switch status {
	case entity.TodoStatusPending:
		return todov1.TodoStatus_TODO_STATUS_PENDING
	case entity.TodoStatusInProgress:
		return todov1.TodoStatus_TODO_STATUS_IN_PROGRESS
	case entity.TodoStatusDone:
		return todov1.TodoStatus_TODO_STATUS_DONE
	default:
		return todov1.TodoStatus_TODO_STATUS_UNSPECIFIED
	}
}
