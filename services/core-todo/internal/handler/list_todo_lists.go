package handler

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/resourcename"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	grpcerrors "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
)

func (h *TodoHandler) ListTodoLists(ctx context.Context, req *todov1.ListTodoListsRequest) (*todov1.ListTodoListsResponse, error) {
	parsed, err := resourcename.ParseUserResourceName(req.GetParent())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid parent: %v", err))
	}

	ownerID := entity.UserID(parsed)

	opts := gatewayinput.ListTodoListsOptions{
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
