// Package middleware implements HTTP middleware for the BFF service, including authentication and other cross-cutting concerns for handling HTTP requests.
package middleware

import (
	"net/http"
	"strings"

	authpkg "github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/infra/jwt"
)

func AuthMiddleware(jwtManager *jwt.JwtManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			payload, err := jwtManager.Verify(tokenStr)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := authpkg.SetUserContext(r.Context(), payload.UserID, payload.Roles)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
