package todo

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/resourcename"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	grpcerrors "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *TodoHandler) DeleteTodo(ctx context.Context, req *todov1.DeleteTodoRequest) (*emptypb.Empty, error) {
	// Parse resource name
	parsed, err := resourcename.ParseTodoResourceName(req.GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid todo name: %v", err))
	}

	// Build input
	in := &input.TodoDeleter{
		ID: entity.TodoID(parsed.TodoID),
	}

	// Execute
	if _, err := h.todoDeleter.Delete(ctx, in); err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	return &emptypb.Empty{}, nil
}
