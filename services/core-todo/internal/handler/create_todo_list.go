package handler

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	grpcerrors "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/helper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
)

func (h *TodoHandler) CreateTodoList(ctx context.Context, req *todov1.CreateTodoListRequest) (*todov1.CreateTodoListResponse, error) {
	requesterID, err := helper.ExtractRequesterID(auth.GetUserID(ctx))
	if err != nil {
		return nil, err
	}

	// Build input
	in := &input.TodoListCreator{
		Name:        req.GetDisplayName(),
		RequesterID: requesterID,
	}

	// Execute
	res, err := h.todoListCreator.Create(ctx, in)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	return &todov1.CreateTodoListResponse{
		TodoList: mapper.TodoListToPb(res.TodoList),
	}, nil
}
