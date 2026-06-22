package gateway

import "context"

type Transactor interface {
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
