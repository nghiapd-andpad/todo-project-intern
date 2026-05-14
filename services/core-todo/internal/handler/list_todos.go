package handler

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	"github.com/nghiapd-andpad/todo-project-intern/pkg/resourcename"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	grpcerrors "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/helper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
)

func (h *TodoHandler) ListTodos(ctx context.Context, req *todov1.ListTodosRequest) (*todov1.ListTodosResponse, error) {
	// Parse parent resource name
	parent, err := resourcename.ParseTodoListResourceName(req.GetParent())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid parent: %v", err))
	}

	requesterID, err := helper.ExtractRequesterID(auth.GetUserID(ctx))
	if err != nil {
		return nil, err
	}

	// Build input
	in := &input.TodoLister{
		TodoListID:  entity.TodoListID(parent.TodoListID),
		RequesterID: requesterID,
		Offset:      int(req.GetOffset()),
		Limit:       int(req.GetPageSize()),
	}

	// Optional filters
	if s := mapper.PbToStatus(req.GetStatusFilter()); s != nil {
		in.Status = s
	}
	if p := mapper.PbToPriority(req.GetPriorityFilter()); p != nil {
		in.Priority = p
	}
	if ts := req.GetTitleSearch(); ts != "" {
		in.TitleSearch = &ts
	}

	// Execute
	res, err := h.todoLister.List(ctx, in)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	// Map response
	pbTodos := make([]*todov1.Todo, len(res.Page.Items))
	for i, t := range res.Page.Items {
		pbTodos[i] = mapper.TodoToPb(t)
	}

	return &todov1.ListTodosResponse{
		Todos: pbTodos,
		Total: res.Page.TotalCount,
	}, nil
}
