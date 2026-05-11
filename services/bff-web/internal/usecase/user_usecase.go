package usecase

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/output"
)

type UserGetter interface {
	GetByID(ctx context.Context, id string) (*output.UserOutput, error)
	GetByIDs(ctx context.Context, ids []string) ([]*output.UserOutput, error)
	GetByUsername(ctx context.Context, username string) (*output.UserOutput, error)
	GetByEmail(ctx context.Context, email string) (*output.UserOutput, error)
}
