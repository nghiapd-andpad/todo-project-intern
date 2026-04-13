package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/user-service/internal/domain/entity"
)

type UserQueriesGateway interface {
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByID(ctx context.Context, id entity.UserID) (*entity.User, error)
}
