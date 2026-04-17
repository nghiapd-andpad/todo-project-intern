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

// Update a Todo by Resource Name: users/{u_id}/todo-lists/{l_id}/todos/{t_id}
func (h *TodoHandler) UpdateTodo(ctx context.Context, req *todov1.UpdateTodoRequest) (*todov1.Todo, error) {
	if req.GetTodo() == nil {
		return nil, status.Error(codes.InvalidArgument, "todo resource is required")
	}

	// Parse ID từ field 'name' của object Todo gửi lên
	todoID, err := parseTodoID(req.GetTodo().GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Build Input
	title := req.GetTodo().GetTitle()
	desc := req.GetTodo().GetDescription()

	in := &input.TodoUpdater{
		ID:          entity.TodoID(todoID),
		Title:       &title,
		Description: &desc,
	}

	// Execute
	out, err := h.TodoUpdater.Update(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update todo: %v", err)
	}

	// Map
	return mapper.TodoToPb(out.Todo), nil
}
