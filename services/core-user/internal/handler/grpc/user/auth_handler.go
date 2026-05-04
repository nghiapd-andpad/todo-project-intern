package user

import (
	"context"

	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	grpcerrors "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/handler/grpc/errors"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/handler/grpc/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user/input"
)

func (h *UserHandler) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	// Build input
	in := &input.UserRegister{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
	}

	// Execute
	out, err := h.userCreator.Register(ctx, in)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	// Map response
	return &userv1.RegisterResponse{
		User: mapper.UserToPb(out.User),
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	// Build input
	in := &input.UserLogin{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	// Execute
	out, err := h.userAuthenticator.Login(ctx, in)
	if err != nil {
		return nil, grpcerrors.ToGRPC(err)
	}

	// Map response
	return &userv1.LoginResponse{
		AccessToken: out.AccessToken,
		User:        mapper.UserToPb(out.User),
	}, nil
}
