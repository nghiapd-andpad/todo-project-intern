package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		// get user ID and roles from metadata
		userIDs := md.Get("x-user-id")
		var userID string
		if len(userIDs) > 0 {
			userID = userIDs[0]
		}

		rolesHeader := md.Get("x-user-roles")
		var roles []string
		if len(rolesHeader) > 0 {
			roles = strings.Split(rolesHeader[0], ",")
		}

		// Inject user ID and roles into context
		newCtx := SetUserContext(ctx, userID, roles)
		return handler(newCtx, req)
	}
}
