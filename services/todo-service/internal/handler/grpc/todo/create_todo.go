package todo

import (
	"context"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/handler/grpc/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/todo-service/internal/usecase/todos/input"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create a new Todo under Parent Resource Name: users/{u_id}/todo-lists/{l_id}
func (h *TodoHandler) CreateTodo(ctx context.Context, req *todov1.CreateTodoRequest) (*todov1.Todo, error) {
	// Validate Parent
	if req.GetParent() == "" {
		return nil, status.Error(codes.InvalidArgument, "parent is required")
	}

	// Build Input
	in := &input.TodoCreator{
		Title:       req.GetTitle(),
		Description: &req.Description,
	}

	// Execute
	out, err := h.TodoCreator.Create(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create todo: %v", err)
	}

	// Map
	return mapper.TodoToPb(out.Todo), nil
}
