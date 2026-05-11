package handler

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/resourcename"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	grpcerrors "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
)

func (h *TodoHandler) UpdateTodo(ctx context.Context, req *todov1.UpdateTodoRequest) (*todov1.UpdateTodoResponse, error) {
	// Validate
	if req.GetTodo() == nil {
		return nil, status.Error(codes.InvalidArgument, "todo is required")
	}
	if req.GetUpdateMask() == nil || len(req.GetUpdateMask().GetPaths()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "update_mask is required")
	}

	// Parse resource name
	parsed, err := resourcename.ParseTodoResourceName(req.GetTodo().GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid todo name: %v", err))
	}

	// Build input
	in := &input.TodoUpdater{
		ID: entity.TodoID(parsed.TodoID),
	}

	for _, path := range req.GetUpdateMask().GetPaths() {
		switch path {
		case "title":
			t := req.GetTodo().GetTitle()
			in.Title = &t
		case "description":
			d := req.GetTodo().GetDescription()
			in.Description = &d
		case "status":
			in.Status = mapper.PbToStatus(req.GetTodo().GetStatus())
		case "priority":
			in.Priority = mapper.PbToPriority(req.GetTodo().GetPriority())
		case "due_date":
			d := req.GetTodo().GetDueDate()
			if d != nil {
				dateStr := d.AsTime().Format("2006-01-02")
				in.DueDate = &dateStr
			}
		case "assignee_id":
			if req.GetTodo().GetAssigneeId() != 0 {
				aID := entity.UserID(req.GetTodo().GetAssigneeId())
				in.AssigneeID = &aID
			}
		}
	}

	// Execute
	out, err := h.todoUpdater.Update(ctx, in)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	// Map response
	return &todov1.UpdateTodoResponse{
		Todo: mapper.TodoToPb(out.Todo),
	}, nil
}
