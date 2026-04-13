package auth

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	secretKey string
}

func NewAuthInterceptor(secretKey string) *AuthInterceptor {
	return &AuthInterceptor{secretKey: secretKey}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Get metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		// Get Bearer Token
		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		accessToken := strings.TrimPrefix(values[0], "Bearer ")

		// Verify JWT
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(i.secretKey), nil
		})

		if err != nil || !token.Valid {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// Get user ID and roles from claims
		userID := int64((*claims)["user_id"].(float64))

		var roles []string
		if r, ok := (*claims)["roles"].([]interface{}); ok {
			for _, v := range r {
				roles = append(roles, v.(string))
			}
		}

		// Set user ID and roles in context
		newCtx := SetUserContext(ctx, userID, roles)
		return handler(newCtx, req)
	}
}
