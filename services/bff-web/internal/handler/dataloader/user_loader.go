package dataloader

import (
	"context"
	"net/http"
	"time"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	userusecase "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/user"
)

type contextKey string

const loadersKey contextKey = "dataloaders"

type Loaders struct {
	UserByID *UserLoader
}

type UserLoader struct {
	userGetter userusecase.UserGetter
	wait       time.Duration
	maxBatch   int
}

func NewUserLoader(userGetter userusecase.UserGetter) *UserLoader {
	return &UserLoader{
		userGetter: userGetter,
		wait:       1 * time.Millisecond,
		maxBatch:   100,
	}
}

func (l *UserLoader) Load(ctx context.Context, id string) (*entity.User, error) {
	return l.userGetter.GetByID(ctx, id)
}

func NewLoaders(userGetter userusecase.UserGetter) *Loaders {
	return &Loaders{
		UserByID: NewUserLoader(userGetter),
	}
}

func Middleware(userGetter userusecase.UserGetter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loaders := NewLoaders(userGetter)
			ctx := context.WithValue(r.Context(), loadersKey, loaders)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func For(ctx context.Context) *Loaders {
	loaders, ok := ctx.Value(loadersKey).(*Loaders)
	if !ok {
		return nil
	}
	return loaders
}
