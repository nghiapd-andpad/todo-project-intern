package directive

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain"
)

func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	if _, ok := auth.GetUserID(ctx); !ok {
		return nil, domain.ErrUnauthorized
	}
	return next(ctx)
}

func HasRoles(ctx context.Context, obj interface{}, next graphql.Resolver, roles []string) (interface{}, error) {
	userID, ok := auth.GetUserID(ctx)
	if !ok || userID == "" {
		return nil, domain.ErrUnauthorized
	}

	userRoles, ok := auth.GetRoles(ctx)
	if !ok {
		return nil, domain.ErrForbidden
	}

	authorized := false
	for _, requiredRole := range roles {
		if contains(userRoles, requiredRole) {
			authorized = true
			break
		}
	}

	if !authorized {
		return nil, domain.ErrForbidden
	}

	return next(ctx)
}

func contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
