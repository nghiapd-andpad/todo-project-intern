// Package dataloader implements data loading and batching for entities, allowing efficient retrieval of user data in the BFF service.
package dataloader

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/graph-gophers/dataloader/v7"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	userusecase "github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/user"
)

type contextKey string

const loadersKey contextKey = "dataloaders"

type Loaders struct {
	UserByID *dataloader.Loader[string, *entity.User]
}

type UserBatcher struct {
	userGetter userusecase.UserGetter
}

func (b *UserBatcher) BatchGetUsers(ctx context.Context, ids []string) []*dataloader.Result[*entity.User] {
	fmt.Printf("fetching bacth for %d IDs: %v\n", len(ids), ids)

	users, err := b.userGetter.GetByIDs(ctx, ids)
	if err != nil {
		// return all errors if the batch query fails
		results := make([]*dataloader.Result[*entity.User], len(ids))
		for i := range results {
			results[i] = &dataloader.Result[*entity.User]{Error: fmt.Errorf("failed to get users: %w", err)}
		}
		return results
	}

	userMap := make(map[string]*entity.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	results := make([]*dataloader.Result[*entity.User], len(ids))
	for i, id := range ids {
		if u, ok := userMap[id]; ok {
			results[i] = &dataloader.Result[*entity.User]{Data: u}
		} else {
			// return error if a specific user is not found
			results[i] = &dataloader.Result[*entity.User]{
				Error: fmt.Errorf("user not found: %s", id),
			}
		}
	}
	return results
}

func NewLoaders(userGetter userusecase.UserGetter, cfg *config.Config) *Loaders {
	batcher := &UserBatcher{userGetter: userGetter}

	return &Loaders{
		UserByID: dataloader.NewBatchedLoader(
			batcher.BatchGetUsers,
			dataloader.WithWait[string, *entity.User](time.Duration(cfg.DataLoaderWait)*time.Millisecond),
			dataloader.WithBatchCapacity[string, *entity.User](cfg.DataLoaderBatchSize),
		),
	}
}

func Middleware(userGetter userusecase.UserGetter, cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loaders := NewLoaders(userGetter, cfg)
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
