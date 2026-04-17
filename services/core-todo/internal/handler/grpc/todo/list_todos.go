package todo

import (
	"context"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Get List of Todos by Parent Resource Name: users/{u_id}/todo-lists/{l_id}
func (h *TodoHandler) ListTodos(ctx context.Context, req *todov1.ListTodosRequest) (*todov1.ListTodosResponse, error) {
	// Validate Parent
	if req.GetParent() == "" {
		return nil, status.Error(codes.InvalidArgument, "parent is required")
	}

	// Build Input
	in := &input.TodoLister{
		Parent: req.GetParent(),
	}

	// Execute
	out, err := h.TodoListReader.List(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list todos: %v", err)
	}

	// Map List
	pbTodos := make([]*todov1.Todo, len(out.Todos))
	for i, t := range out.Todos {
		pbTodos[i] = mapper.TodoToPb(t)
	}

	return &todov1.ListTodosResponse{
		Todos: pbTodos,
	}, nil
}
