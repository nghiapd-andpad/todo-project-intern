package errors

import (
	"errors"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToStatus(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := status.FromError(err); ok {
		return err
	}

	switch {
	case errors.Is(err, entity.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, entity.ErrUsernameAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, entity.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, entity.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, entity.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, "invalid credentials")

	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
