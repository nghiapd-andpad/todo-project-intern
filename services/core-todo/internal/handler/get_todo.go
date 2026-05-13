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

func (h *TodoHandler) GetTodo(ctx context.Context, req *todov1.GetTodoRequest) (*todov1.GetTodoResponse, error) {
	// Parse resource name
	parsed, err := resourcename.ParseTodoResourceName(req.GetName())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid todo name: %v", err)
	}

	requesterID, err := helper.ExtractRequesterID(auth.GetUserID(ctx))
	if err != nil {
		return nil, err
	}

	// Build input
	in := &input.TodoGetter{
		TodoID:      entity.TodoID(parsed.TodoID),
		TodoListID:  entity.TodoListID(parsed.TodoListID),
		RequesterID: requesterID,
	}

	// Execute
	res, err := h.todoGetter.Get(ctx, in)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	// Map response
	return &todov1.GetTodoResponse{
		Todo: mapper.TodoToPb(res.Todo),
	}, nil
}
