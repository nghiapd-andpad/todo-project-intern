// Package mapper provides functions to convert between protobuf messages and domain entities for the core service.
package mapper

import (
	"fmt"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
)

func TodoListFromPb(pb *todov1.TodoList) *entity.TodoList {
	if pb == nil {
		return nil
	}
	return &entity.TodoList{
		Name:        pb.Name,
		DisplayName: pb.DisplayName,
		CreatedAt:   pb.CreatedAt.AsTime(),
		UpdatedAt:   pb.UpdatedAt.AsTime(),
	}
}

func TodoFromPb(pb *todov1.Todo) *entity.Todo {
	if pb == nil {
		return nil
	}
	t := &entity.Todo{
		Name:      pb.Name,
		Title:     pb.Title,
		Status:    TodoStatusFromPb(pb.Status),
		Priority:  PriorityFromPb(pb.Priority),
		CreatorID: fmt.Sprintf("users/%d", pb.CreatorId),
		CreatedAt: pb.CreatedAt.AsTime(),
		UpdatedAt: pb.UpdatedAt.AsTime(),
	}
	if pb.Description != "" {
		t.Description = &pb.Description
	}
	if pb.DueDate != nil {
		dt := pb.DueDate.AsTime()
		t.DueDate = &dt
	}
	if pb.AssigneeId != 0 {
		aID := fmt.Sprintf("users/%d", pb.AssigneeId)
		t.AssigneeID = &aID
	}
	return t
}

func TodoStatusFromPb(s todov1.TodoStatus) entity.TodoStatus {
	switch s {
	case todov1.TodoStatus_TODO_STATUS_PENDING:
		return entity.TodoStatusPending
	case todov1.TodoStatus_TODO_STATUS_IN_PROGRESS:
		return entity.TodoStatusInProgress
	case todov1.TodoStatus_TODO_STATUS_DONE:
		return entity.TodoStatusDone
	default:
		return entity.TodoStatusPending
	}
}

func TodoStatusToPb(s entity.TodoStatus) todov1.TodoStatus {
	switch s {
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

func PriorityFromPb(p todov1.Priority) entity.Priority {
	switch p {
	case todov1.Priority_PRIORITY_LOW:
		return entity.PriorityLow
	case todov1.Priority_PRIORITY_MEDIUM:
		return entity.PriorityMedium
	case todov1.Priority_PRIORITY_HIGH:
		return entity.PriorityHigh
	case todov1.Priority_PRIORITY_URGENT:
		return entity.PriorityUrgent
	default:
		return entity.PriorityMedium
	}
}

func PriorityToPb(p entity.Priority) todov1.Priority {
	switch p {
	case entity.PriorityLow:
		return todov1.Priority_PRIORITY_LOW
	case entity.PriorityMedium:
		return todov1.Priority_PRIORITY_MEDIUM
	case entity.PriorityHigh:
		return todov1.Priority_PRIORITY_HIGH
	case entity.PriorityUrgent:
		return todov1.Priority_PRIORITY_URGENT
	default:
		return todov1.Priority_PRIORITY_UNSPECIFIED
	}
}
