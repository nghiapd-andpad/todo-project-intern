//go:generate mockgen -destination=mock/transactor_mock.go -source=transactor.go -package mock
package gateway

import "context"

type Transactor interface {
	Transaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
