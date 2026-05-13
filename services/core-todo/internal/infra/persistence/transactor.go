package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
)

type txKey struct{}

type Transactor struct {
	db *gorm.DB
}

func NewTransactor(db *gorm.DB) *Transactor {
	return &Transactor{db: db}
}

var _ gateway.Transactor = (*Transactor)(nil)

func (t *Transactor) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		if err := fn(txCtx); err != nil {
			return err // ROLLBACK
		}
		return nil // COMMIT
	})
}

func connFromContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return db.WithContext(ctx)
}
