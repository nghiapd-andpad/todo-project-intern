package todo

import (
	"context"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Get a Todo by Resource Name: users/{u_id}/todo-lists/{l_id}/todos/{t_id}
func (h *TodoHandler) GetTodo(ctx context.Context, req *todov1.GetTodoRequest) (*todov1.Todo, error) {
	// Parse
	todoID, err := parseTodoID(req.GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Build Input
	in := &input.TodoGetter{
		ID: entity.TodoID(todoID),
	}

	// Execute
	out, err := h.TodoGetter.Get(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get todo: %v", err)
	}

	// Map
	return mapper.TodoToPb(out.Todo), nil
}
