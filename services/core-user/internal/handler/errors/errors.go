// Package errors provides utilities for converting application errors to gRPC status errors.
package errors

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
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
		case entity.ErrInvalidCredentials:
			return status.Error(codes.Unauthenticated, appErr.Message)
		case entity.ErrUsernameAlreadyExists:
			return status.Error(codes.AlreadyExists, appErr.Message)
		case entity.ErrEmailAlreadyExists:
			return status.Error(codes.AlreadyExists, appErr.Message)
		case entity.ErrInvalidToken:
			return status.Error(codes.Unauthenticated, appErr.Message)
		case entity.ErrExpiredToken:
			return status.Error(codes.Unauthenticated, appErr.Message)
		case entity.ErrInternal:
			return status.Error(codes.Internal, appErr.Message)
		}
	}

	return status.Error(codes.Internal, "internal server error")
}
