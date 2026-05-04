package todo

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/resourcename"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	grpcerrors "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
)

func (h *TodoHandler) UpdateTodoList(ctx context.Context, req *todov1.UpdateTodoListRequest) (*todov1.UpdateTodoListResponse, error) {
	if req.GetTodoList() == nil {
		return nil, status.Error(codes.InvalidArgument, "todo_list is required")
	}
	if req.GetUpdateMask() == nil || len(req.GetUpdateMask().GetPaths()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "update_mask is required")
	}

	parsed, err := resourcename.ParseTodoListResourceName(req.GetTodoList().GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid todo list name: %v", err))
	}

	in := &input.TodoListUpdater{
		ID: entity.TodoListID(parsed.TodoListID),
	}

	for _, path := range req.GetUpdateMask().GetPaths() {
		switch path {
		case "display_name":
			n := req.GetTodoList().GetDisplayName()
			in.Name = &n
		}
	}

	out, err := h.todoListUpdater.Update(ctx, in)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	return &todov1.UpdateTodoListResponse{
		TodoList: mapper.TodoListToPb(out.TodoList),
	}, nil
}
