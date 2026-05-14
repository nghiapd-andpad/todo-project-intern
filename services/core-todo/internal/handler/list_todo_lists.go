package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	"github.com/nghiapd-andpad/todo-project-intern/pkg/resourcename"
	todov1 "github.com/nghiapd-andpad/todo-project-intern/proto/todo/v1"
	grpcerrors "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/helper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/handler/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/usecase/input"
)

func (h *TodoHandler) ListTodoLists(ctx context.Context, req *todov1.ListTodoListsRequest) (*todov1.ListTodoListsResponse, error) {
	if _, err := resourcename.ParseUserResourceName(req.GetParent()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
	}

	// requesterID from auth metadata is the source of truth for authorization.
	requesterID, err := helper.ExtractRequesterID(auth.GetUserID(ctx))
	if err != nil {
		return nil, err
	}

	// Build input
	in := &input.TodoListLister{
		RequesterID: requesterID,
		Filter:      input.TodoListFilterAll,
		Offset:      int(req.GetOffset()),
		Limit:       int(req.GetPageSize()),
	}
	if ns := req.GetNameSearch(); ns != "" {
		in.NameSearch = &ns
	}

	// Execute
	res, err := h.todoListLister.List(ctx, in)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	pbLists := make([]*todov1.TodoList, len(res.Page.Items))
	for i, tl := range res.Page.Items {
		pbLists[i] = mapper.TodoListToPb(tl)
	}

	return &todov1.ListTodoListsResponse{
		TodoLists: pbLists,
		Total:     res.Page.TotalCount,
	}, nil
}
