package todo

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	grpcerrors "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/grpc/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/todos/input"
)

func (h *TodoHandler) ListTodoLists(ctx context.Context, req *todov1.ListTodoListsRequest) (*todov1.ListTodoListsResponse, error) {
	// Parse "users/{user_id}"
	var userID int64
	if _, err := fmt.Sscanf(req.GetParent(), "users/%d", &userID); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid parent format: %s", req.GetParent()))
	}

	ownerID := entity.UserID(userID)
	opts := gateway.ListTodoListsOptions{
		OwnerID: &ownerID,
		Offset:  int(req.GetOffset()),
		Limit:   int(req.GetPageSize()),
	}
	if ns := req.GetNameSearch(); ns != "" {
		opts.NameSearch = &ns
	}

	out, err := h.todoListLister.List(ctx, &input.TodoListLister{Opts: opts})
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	pbLists := make([]*todov1.TodoList, len(out.TodoLists))
	for i, tl := range out.TodoLists {
		pbLists[i] = mapper.TodoListToPb(tl)
	}

	return &todov1.ListTodoListsResponse{
		TodoLists: pbLists,
		Total:     out.Total,
	}, nil
}
