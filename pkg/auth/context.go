package auth

import (
	"context"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
	rolesKey  contextKey = "roles"
)

// SetUserContext sets user ID and roles in context.
func SetUserContext(ctx context.Context, userID string, roles []string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	return context.WithValue(ctx, rolesKey, roles)
}

// GetUserID gets user ID from context.
func GetUserID(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(userIDKey).(string)
	return uid, ok
}

// GetRoles gets roles from context.
func GetRoles(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(rolesKey).([]string)
	return roles, ok
}
