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

func (h *TodoHandler) GetTodoList(ctx context.Context, req *todov1.GetTodoListRequest) (*todov1.GetTodoListResponse, error) {
	parsed, err := resourcename.ParseTodoListResourceName(req.GetName())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid todo list name: %v", err)
	}

	requesterID, err := helper.ExtractRequesterID(auth.GetUserID(ctx))
	if err != nil {
		return nil, err
	}

	// Build input
	in := &input.TodoListGetter{
		TodoListID:  entity.TodoListID(parsed.TodoListID),
		RequesterID: requesterID,
	}

	// Execute
	res, err := h.todoListGetter.Get(ctx, in)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	return &todov1.GetTodoListResponse{TodoList: mapper.TodoListToPb(res.TodoList)}, nil
}
