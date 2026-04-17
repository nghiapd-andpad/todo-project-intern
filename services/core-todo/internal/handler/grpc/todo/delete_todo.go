package todo

import (
	"context"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Delete a Todo by Resource Name: users/{u_id}/todo-lists/{l_id}/todos/{t_id}
func (h *TodoHandler) DeleteTodo(ctx context.Context, req *todov1.DeleteTodoRequest) (*emptypb.Empty, error) {
	// Parse
	todoID, err := parseTodoID(req.GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Build Input
	in := &input.TodoDeleter{
		ID: entity.TodoID(todoID),
	}

	// Execute
	_, err = h.TodoDeleter.Delete(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete todo: %v", err)
	}

	return &emptypb.Empty{}, nil
}
