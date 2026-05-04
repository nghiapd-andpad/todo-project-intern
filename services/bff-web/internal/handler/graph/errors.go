// Package graph implements the GraphQL server for the BFF service, including schema definitions, resolvers, and error handling for GraphQL operations.
package graph

import (
	"context"
	"errors"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
)

func ErrorPresenter(ctx context.Context, err error) *gqlerror.Error {
	var appErr *entity.AppError
	if errors.As(err, &appErr) {
		return &gqlerror.Error{
			Message: appErr.Message,
			Extensions: map[string]interface{}{
				"code": string(appErr.Code),
			},
		}
	}

	// gRPC status error
	if st, ok := status.FromError(err); ok && st.Code() != codes.OK {
		code := grpcCodeToExtension(st.Code())
		return &gqlerror.Error{
			Message: st.Message(),
			Extensions: map[string]interface{}{
				"code": code,
			},
		}
	}

	// Unwrap gRPC error
	unwrapped := err
	for unwrapped != nil {
		if st, ok := status.FromError(unwrapped); ok && st.Code() != codes.OK {
			return &gqlerror.Error{
				Message: st.Message(),
				Extensions: map[string]interface{}{
					"code": grpcCodeToExtension(st.Code()),
				},
			}
		}
		unwrapped = errors.Unwrap(unwrapped)
	}

	return &gqlerror.Error{
		Message: "internal server error",
		Extensions: map[string]interface{}{
			"code": "INTERNAL",
		},
	}
}

func grpcCodeToExtension(code codes.Code) string {
	switch code {
	case codes.NotFound:
		return "NOT_FOUND"
	case codes.AlreadyExists:
		return "ALREADY_EXISTS"
	case codes.Unauthenticated:
		return "UNAUTHENTICATED"
	case codes.PermissionDenied:
		return "FORBIDDEN"
	case codes.InvalidArgument:
		return "INVALID_ARGUMENT"
	case codes.Internal:
		return "INTERNAL"
	default:
		return "INTERNAL"
	}
}
