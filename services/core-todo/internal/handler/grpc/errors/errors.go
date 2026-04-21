package errors

import (
	"errors"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPC(err error) error {
	if err == nil {
		return nil
	}

	var appErr *entity.AppError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case entity.ErrNotFound:
			return status.Error(codes.NotFound, appErr.Message)
		case entity.ErrInvalidParameter:
			return status.Error(codes.InvalidArgument, appErr.Message)
		case entity.ErrAuthZ:
			return status.Error(codes.PermissionDenied, appErr.Message)
		case entity.ErrAuthN:
			return status.Error(codes.Unauthenticated, appErr.Message)
		case entity.ErrInternal:
			return status.Error(codes.Internal, appErr.Message)
		}
	}

	// Unexpected error
	return status.Error(codes.Internal, "internal server error")
}
