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

func (h *TodoHandler) UpdateTodo(ctx context.Context, req *todov1.UpdateTodoRequest) (*todov1.UpdateTodoResponse, error) {
	if req.GetUpdateMask() == nil || len(req.GetUpdateMask().GetPaths()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "update_mask is required")
	}

	parsed, err := resourcename.ParseTodoResourceName(req.GetTodo().GetName())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid todo name: %v", err)
	}

	requesterID, err := helper.ExtractRequesterID(auth.GetUserID(ctx))
	if err != nil {
		return nil, err
	}

	// Build input
	in := &input.TodoUpdater{
		TodoID:      entity.TodoID(parsed.TodoID),
		TodoListID:  entity.TodoListID(parsed.TodoListID),
		RequesterID: requesterID,
	}

	for _, path := range req.GetUpdateMask().GetPaths() {
		switch path {
		case "title":
			t := req.GetTodo().GetTitle()
			in.Fields.Title = &t
		case "description":
			d := req.GetTodo().GetDescription()
			in.Fields.Description = &d
		case "status":
			in.Fields.Status = mapper.PbToStatus(req.GetTodo().GetStatus())
		case "priority":
			in.Fields.Priority = mapper.PbToPriority(req.GetTodo().GetPriority())
		case "due_date":
			if d := req.GetTodo().GetDueDate(); d != nil {
				t := d.AsTime()
				in.Fields.DueDate = &t
			}
		case "assignee_id":
			if id := req.GetTodo().GetAssigneeId(); id != 0 {
				aID := entity.UserID(id)
				in.Fields.AssigneeID = &aID
			}
		}
	}

	// Execute
	res, err := h.todoUpdater.Update(ctx, in)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	// Map response
	return &todov1.UpdateTodoResponse{
		Todo: mapper.TodoToPb(res.Todo),
	}, nil
}
