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
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
)

func (h *TodoHandler) DeleteTodoList(ctx context.Context, req *todov1.DeleteTodoListRequest) (*todov1.DeleteTodoListResponse, error) {
	parsed, err := resourcename.ParseTodoListResourceName(req.GetName())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid todo list name: %v", err)
	}

	requesterID, err := helper.ExtractRequesterID(auth.GetUserID(ctx))
	if err != nil {
		return nil, err
	}

	// Build input
	in := &input.TodoListDeleter{
		TodoListID:  entity.TodoListID(parsed.TodoListID),
		RequesterID: requesterID,
	}

	if _, err := h.todoListDeleter.Delete(ctx, in); err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	return &todov1.DeleteTodoListResponse{}, nil
}
