package gateway

import (
	"context"
	"time"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
)

type TokenPayload struct {
	UserID entity.UserID
	Roles  []string
}

type TokenManager interface {
	Generate(ctx context.Context, payload TokenPayload, duration time.Duration) (string, error)
}
