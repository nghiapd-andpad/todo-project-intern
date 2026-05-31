package input

import (
	"time"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

type CreateIdempotencyRecord struct {
	UserID         entity.UserID
	Operation      string
	IdempotencyKey string
	ExpiresAt      time.Time
}
