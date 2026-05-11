package handler

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	grpcerrors "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
)

func (h *TodoHandler) CreateTodoList(ctx context.Context, req *todov1.CreateTodoListRequest) (*todov1.CreateTodoListResponse, error) {
	// Extract OwnerID from auth context
	userIDStr, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing user id in context")
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("invalid user id: %v", err))
	}

	out, err := h.todoListCreator.Create(ctx, &input.TodoListCreator{
		Name:    req.GetDisplayName(),
		OwnerID: entity.UserID(userID),
	})
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	return &todov1.CreateTodoListResponse{
		TodoList: mapper.TodoListToPb(out.TodoList),
	}, nil
}
