package output

import (
	"time"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

type IdempotencyStatus string

const (
	IdempotencyStatusProcessing IdempotencyStatus = "PROCESSING"
	IdempotencyStatusCompleted  IdempotencyStatus = "COMPLETED"
	IdempotencyStatusFailed     IdempotencyStatus = "FAILED"
)

type IdempotencyRecord struct {
	ID             int64
	UserID         entity.UserID
	Operation      string
	IdempotencyKey string
	Status         IdempotencyStatus
	ResourceType   *string
	ResourceID     *int64
	ExpiresAt      time.Time
}
