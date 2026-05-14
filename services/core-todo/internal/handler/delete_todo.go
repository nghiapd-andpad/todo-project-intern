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

func (h *TodoHandler) DeleteTodo(ctx context.Context, req *todov1.DeleteTodoRequest) (*todov1.DeleteTodoResponse, error) {
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
	in := &input.TodoDeleter{
		TodoID:      entity.TodoID(parsed.TodoID),
		TodoListID:  entity.TodoListID(parsed.TodoListID),
		RequesterID: requesterID,
	}

	// Execute
	if _, err := h.todoDeleter.Delete(ctx, in); err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	return &todov1.DeleteTodoResponse{}, nil
}
