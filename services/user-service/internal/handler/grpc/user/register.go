package user

import (
	"context"

	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/user-service/internal/handler/grpc/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/user-service/internal/usecase/user/input"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *UserHandler) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	// Build Input
	in := &input.UserRegister{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
	}

	// Execute UseCase
	out, err := h.UserCreator.Register(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register: %v", err)
	}

	// Map to Proto Response
	return &userv1.RegisterResponse{
		User: mapper.UserToPb(out.User),
	}, nil
}
