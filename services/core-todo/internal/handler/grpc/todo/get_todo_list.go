package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/resourcename"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	grpcerrors "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *TodoHandler) GetTodoList(ctx context.Context, req *todov1.GetTodoListRequest) (*todov1.TodoList, error) {
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

	return mapper.TodoListToPb(out.TodoList), nil
}
