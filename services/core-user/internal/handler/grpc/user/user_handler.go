package user

import (
	"context"
	"strconv"

	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	grpcerrors "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/handler/grpc/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/handler/grpc/mapper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *UserHandler) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.UserResponse, error) {
	// Parse user ID
	id, err := strconv.ParseInt(req.GetId(), 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id format")
	}

	out, err := h.userGetter.GetByID(ctx, entity.UserID(id))
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	return &userv1.UserResponse{
		User: mapper.UserToPb(out),
	}, nil
}

func (h *UserHandler) GetUserByUsername(ctx context.Context, req *userv1.GetUserByUsernameRequest) (*userv1.UserResponse, error) {
	if req.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	out, err := h.userGetter.GetByUsername(ctx, req.GetUsername())
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	return &userv1.UserResponse{
		User: mapper.UserToPb(out),
	}, nil
}

func (h *UserHandler) GetUserByEmail(ctx context.Context, req *userv1.GetUserByEmailRequest) (*userv1.UserResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	out, err := h.userGetter.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	return &userv1.UserResponse{
		User: mapper.UserToPb(out),
	}, nil
}

func (h *UserHandler) BatchGetUsers(ctx context.Context, req *userv1.BatchGetUsersRequest) (*userv1.BatchGetUsersResponse, error) {
	if len(req.GetIds()) == 0 {
		return &userv1.BatchGetUsersResponse{}, nil
	}

	// Convert IDs to list of entity.UserID
	userIDs := make([]entity.UserID, 0, len(req.GetIds()))
	for _, idStr := range req.GetIds() {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid user id format: %s", idStr)
		}
		userIDs = append(userIDs, entity.UserID(id))
	}

	// Execute
	dtos, err := h.userGetter.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	// Map response
	pbUsers := make([]*userv1.User, len(dtos))
	for i, out := range dtos {
		pbUsers[i] = mapper.UserToPb(out)
	}

	return &userv1.BatchGetUsersResponse{
		Users: pbUsers,
	}, nil
}
