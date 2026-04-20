package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// get userID and roles from context
		userID, ok := GetUserID(ctx)
		if ok && userID != "" {
			// get metadata
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				md = metadata.New(make(map[string]string))
			} else {
				md = md.Copy()
			}

			// inject userID and roles into metadata
			md.Set("x-user-id", userID)

			if roles, ok := GetRoles(ctx); ok && len(roles) > 0 {
				md.Set("x-user-roles", strings.Join(roles, ","))
			}

			// 4. Inject new metadata into context
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
