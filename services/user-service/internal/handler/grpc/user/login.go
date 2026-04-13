package user

import (
	"context"

	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/user-service/internal/handler/grpc/mapper"
	"github.com/nghiapd-andpad/todo-project-intern/services/user-service/internal/usecase/user/input"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *UserHandler) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	// Build Input
	in := &input.UserLogin{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	// Execute UseCase
	out, err := h.UserAuthenticator.Login(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid username or password")
	}

	// Map to Proto Response
	return &userv1.LoginResponse{
		AccessToken: out.AccessToken,
		User:        mapper.UserToPb(out.User),
	}, nil
}
