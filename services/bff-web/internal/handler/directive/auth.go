// Package directive implements GraphQL directives for the BFF service, including authentication checks and other middleware-like functionality for GraphQL resolvers.
package directive

import (
	"context"

	"github.com/99designs/gqlgen/graphql"

	authpkg "github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
)

// check user authenticated? (@auth)
func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	userID, ok := authpkg.GetUserID(ctx)
	if !ok || userID == "" {
		return nil, entity.NewAuthN("authentication required")
	}
	return next(ctx)
}
