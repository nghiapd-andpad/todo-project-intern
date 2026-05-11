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

func (h *TodoHandler) GetTodoList(ctx context.Context, req *todov1.GetTodoListRequest) (*todov1.GetTodoListResponse, error) {
	parsed, err := resourcename.ParseTodoListResourceName(req.GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid todo list name: %v", err))
	}

	out, err := h.todoListGetter.Get(ctx, &input.TodoListGetter{
		ID: entity.TodoListID(parsed.TodoListID),
	})
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	return &todov1.GetTodoListResponse{
		TodoList: mapper.TodoListToPb(out.TodoList),
	}, nil
}
