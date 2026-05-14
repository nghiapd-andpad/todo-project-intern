package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	"github.com/nghiapd-andpad/todo-project-intern/pkg/resourcename"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	grpcerrors "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/helper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
)

func (h *TodoHandler) CreateTodo(ctx context.Context, req *todov1.CreateTodoRequest) (*todov1.CreateTodoResponse, error) {
	// Parse parent resource name
	parent, err := resourcename.ParseTodoListResourceName(req.GetParent())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
	}

	// Extract creator from auth context
	requesterID, err := helper.ExtractRequesterID(auth.GetUserID(ctx))
	if err != nil {
		return nil, err
	}

	dueDate, err := helper.ParseDueDate(req.GetDueDate())
	if err != nil {
		return nil, err
	}

	// Build input
	in := &input.TodoCreator{
		TodoListID:  entity.TodoListID(parent.TodoListID),
		RequesterID: requesterID,
		Title:       req.GetTitle(),
		Priority:    mapper.PbToPriorityValue(req.GetPriority()),
		DueDate:     dueDate,
	}

	// Optional fields
	if desc := req.GetDescription(); desc != "" {
		in.Description = &desc
	}
	if req.GetAssigneeId() != 0 {
		aID := entity.UserID(req.GetAssigneeId())
		in.AssigneeID = &aID
	}

	// Execute
	res, err := h.todoCreator.Create(ctx, in)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	// Map response
	return &todov1.CreateTodoResponse{
		Todo: mapper.TodoToPb(res.Todo),
	}, nil
}
