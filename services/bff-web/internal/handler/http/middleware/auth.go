package middleware

import (
	"errors"
	"net/http"
	"strings"

	authPkg "github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/auth"
)

func AuthMiddleware(au auth.AuthUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, roles, err := au.Authenticate(r.Context(), token)

			if err != nil {
				if errors.Is(err, domain.ErrUnauthorized) {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				if errors.Is(err, domain.ErrForbidden) {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			ctx := authPkg.SetUserContext(r.Context(), userID, roles)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
