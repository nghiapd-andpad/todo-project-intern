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

func (h *TodoHandler) GetTodo(ctx context.Context, req *todov1.GetTodoRequest) (*todov1.GetTodoResponse, error) {
	// Parse resource name
	parsed, err := resourcename.ParseTodoResourceName(req.GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid todo name: %v", err))
	}

	// Build input
	in := &input.TodoGetter{
		ID: entity.TodoID(parsed.TodoID),
	}

	// Execute
	out, err := h.todoGetter.Get(ctx, in)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	// Map response
	return &todov1.GetTodoResponse{
		Todo: mapper.TodoToPb(out.Todo),
	}, nil
}
