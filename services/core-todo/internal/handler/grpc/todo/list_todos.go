package todo

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/resourcename"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	grpcerrors "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
)

func (h *TodoHandler) ListTodos(ctx context.Context, req *todov1.ListTodosRequest) (*todov1.ListTodosResponse, error) {
	// Parse parent resource name
	parent, err := resourcename.ParseTodoListResourceName(req.GetParent())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid parent: %v", err))
	}

	// Build opts filters
	listID := entity.TodoListID(parent.TodoListID)
	opts := gatewayinput.ListTodosOptions{
		TodoListID: &listID,
		Offset:     int(req.GetOffset()),
		Limit:      int(req.GetPageSize()),
	}

	// Optional filters
	if s := mapper.PbToStatus(req.GetStatusFilter()); s != nil {
		opts.Status = s
	}
	if p := mapper.PbToPriority(req.GetPriorityFilter()); p != nil {
		opts.Priority = p
	}
	if ts := req.GetTitleSearch(); ts != "" {
		opts.TitleSearch = &ts
	}

	// Build input
	in := &input.TodoLister{Opts: opts}

	// Execute
	out, err := h.todoLister.List(ctx, in)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	// Map response
	pbTodos := make([]*todov1.Todo, len(out.Todos))
	for i, t := range out.Todos {
		pbTodos[i] = mapper.TodoToPb(t)
	}

	return &todov1.ListTodosResponse{
		Todos: pbTodos,
		Total: out.Total,
	}, nil
}
