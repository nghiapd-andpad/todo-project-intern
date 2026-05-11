// Package dataloader provides GraphQL dataloaders for batching and caching database requests.
package dataloader

import (
	"context"
	"net/http"

	"github.com/graph-gophers/dataloader/v7"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/output"
)

type UserByIDLoader = dataloader.Loader[string, *output.UserOutput]

type Loaders struct {
	UserByID *UserByIDLoader
}

type contextKey string

const loadersKey contextKey = "dataloaders"

type UserFetcher interface {
	GetByIDs(ctx context.Context, ids []string) ([]*output.UserOutput, error)
}

// Middleware injects per-request dataloaders into context.
func Middleware(fetcher UserFetcher, cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loaders := &Loaders{
				UserByID: newUserByIDLoader(r.Context(), fetcher),
			}
			ctx := context.WithValue(r.Context(), loadersKey, loaders)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// For retrieves loaders from context.
func For(ctx context.Context) *Loaders {
	loaders, _ := ctx.Value(loadersKey).(*Loaders)
	return loaders
}

func newUserByIDLoader(ctx context.Context, fetcher UserFetcher) *UserByIDLoader {
	batchFn := func(ctx context.Context, ids []string) []*dataloader.Result[*output.UserOutput] {
		users, err := fetcher.GetByIDs(ctx, ids)
		results := make([]*dataloader.Result[*output.UserOutput], len(ids))

		if err != nil {
			for i := range ids {
				results[i] = &dataloader.Result[*output.UserOutput]{Error: err}
			}
			return results
		}

		byID := make(map[string]*output.UserOutput, len(users))
		for _, u := range users {
			byID[u.ID] = u
		}

		for i, id := range ids {
			if u, ok := byID[id]; ok {
				results[i] = &dataloader.Result[*output.UserOutput]{Data: u}
			} else {
				results[i] = &dataloader.Result[*output.UserOutput]{Data: nil}
			}
		}
		return results
	}

	return dataloader.NewBatchedLoader(batchFn)
}
