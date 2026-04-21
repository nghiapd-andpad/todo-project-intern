package todo

import (
	"context"
	"fmt"
	"strconv"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	"github.com/nghiapd-andpad/todo-project-intern/pkg/resourcename"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	grpcerrors "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *TodoHandler) CreateTodo(ctx context.Context, req *todov1.CreateTodoRequest) (*todov1.Todo, error) {
	// Parse parent resource name
	parent, err := resourcename.ParseTodoListResourceName(req.GetParent())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid parent: %v", err))
	}

	// Extract creator from auth context
	userIDStr, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing user id in context")
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid user id in context")
	}

	// Build input
	in := &input.TodoCreator{
		TodoListID: entity.TodoListID(parent.TodoListID),
		Title:      req.GetTitle(),
		Priority:   mapper.PbToPriorityValue(req.GetPriority()), // default Medium if UNSPECIFIED
		CreatorID:  entity.UserID(userID),
	}

	// Optional fields
	if desc := req.GetDescription(); desc != "" {
		in.Description = &desc
	}
	if req.GetDueDate() != "" {
		d := req.GetDueDate()
		in.DueDate = &d
	}
	if req.GetAssigneeId() != 0 {
		aID := entity.UserID(req.GetAssigneeId())
		in.AssigneeID = &aID
	}

	// Execute
	out, err := h.todoCreator.Create(ctx, in)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	// Map response
	return mapper.TodoToPb(out.Todo), nil
}
