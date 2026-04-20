package graph

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"google.golang.org/grpc/status"
)

func PresentError(ctx context.Context, err error) *gqlerror.Error {
	var gqlErr *gqlerror.Error
	originalErr := err
	if errors.As(err, &gqlErr) {
		if gqlErr.Unwrap() != nil {
			originalErr = gqlErr.Unwrap()
		}
	}

	st, ok := status.FromError(originalErr)
	if ok {
		return &gqlerror.Error{
			Message: st.Message(),
			Path:    graphql.GetPath(ctx),
			Extensions: map[string]interface{}{
				"code": st.Code().String(),
			},
		}
	}

	if gqlErr != nil {
		return gqlErr
	}

	return gqlerror.Errorf(err.Error())
}
