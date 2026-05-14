// Package mapper provides functions to convert between domain entities and gRPC protobuf messages.
package mapper

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/resourcename"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

func TodoToPb(ent *entity.Todo) *todov1.Todo {
	if ent == nil {
		return nil
	}

	pb := &todov1.Todo{
		Name: resourcename.TodoResourceName{
			TodoListID: int64(ent.TodoListID),
			TodoID:     int64(ent.ID),
		}.String(),
		Title:     ent.Title,
		Status:    TodoStatusToPb(ent.Status),
		Priority:  PriorityToPb(ent.Priority),
		CreatedAt: timestamppb.New(ent.CreatedAt),
		UpdatedAt: timestamppb.New(ent.UpdatedAt),
	}

	// Optional fields
	if ent.Description != nil {
		pb.Description = *ent.Description
	}
	if ent.DueDate != nil {
		pb.DueDate = timestamppb.New(*ent.DueDate)
	}
	if ent.AssigneeID != nil {
		pb.AssigneeId = int64(*ent.AssigneeID)
	}

	return pb
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

func PbToStatus(s todov1.TodoStatus) *entity.TodoStatus {
	var result entity.TodoStatus
	switch s {
	case todov1.TodoStatus_TODO_STATUS_PENDING:
		result = entity.TodoStatusPending
	case todov1.TodoStatus_TODO_STATUS_IN_PROGRESS:
		result = entity.TodoStatusInProgress
	case todov1.TodoStatus_TODO_STATUS_DONE:
		result = entity.TodoStatusDone
	default:
		return nil // UNSPECIFIED
	}
	return &result
}

func PbToPriority(p todov1.Priority) *entity.Priority {
	var result entity.Priority
	switch p {
	case todov1.Priority_PRIORITY_LOW:
		result = entity.PriorityLow
	case todov1.Priority_PRIORITY_MEDIUM:
		result = entity.PriorityMedium
	case todov1.Priority_PRIORITY_HIGH:
		result = entity.PriorityHigh
	case todov1.Priority_PRIORITY_URGENT:
		result = entity.PriorityUrgent
	default:
		return nil
	}
	return &result
}

func PbToPriorityValue(p todov1.Priority) entity.Priority {
	switch p {
	case todov1.Priority_PRIORITY_LOW:
		return entity.PriorityLow
	case todov1.Priority_PRIORITY_HIGH:
		return entity.PriorityHigh
	case todov1.Priority_PRIORITY_URGENT:
		return entity.PriorityUrgent
	default:
		return entity.PriorityMedium
	}
}
