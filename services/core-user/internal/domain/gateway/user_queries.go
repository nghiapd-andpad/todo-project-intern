package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
)

type UserQueriesGateway interface {
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, id entity.UserID) (*entity.User, error)
	GetByIDs(ctx context.Context, ids []entity.UserID) ([]*entity.User, error)
}
