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
