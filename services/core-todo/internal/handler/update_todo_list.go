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

func (h *TodoHandler) UpdateTodoList(ctx context.Context, req *todov1.UpdateTodoListRequest) (*todov1.UpdateTodoListResponse, error) {
	if req.GetUpdateMask() == nil || len(req.GetUpdateMask().GetPaths()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "update_mask is required")
	}

	parsed, err := resourcename.ParseTodoListResourceName(req.GetTodoList().GetName())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid todo list name: %v", err)
	}

	requesterID, err := helper.ExtractRequesterID(auth.GetUserID(ctx))
	if err != nil {
		return nil, err
	}

	// Build input
	in := &input.TodoListUpdater{
		TodoListID:  entity.TodoListID(parsed.TodoListID),
		RequesterID: requesterID,
	}
	for _, path := range req.GetUpdateMask().GetPaths() {
		if path == "display_name" {
			n := req.GetTodoList().GetDisplayName()
			in.Name = &n
		}
	}

	// Execute
	res, err := h.todoListUpdater.Update(ctx, in)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	return &todov1.UpdateTodoListResponse{TodoList: mapper.TodoListToPb(res.TodoList)}, nil
}
