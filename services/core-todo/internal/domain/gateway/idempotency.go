package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/output"
)

type IdempotencyGateway interface {
	Find(ctx context.Context, userID entity.UserID, operation string, key string) (*output.IdempotencyRecord, error)
	CreateProcessing(ctx context.Context, in *input.CreateIdempotencyRecord) (*output.IdempotencyRecord, error)
	MarkCompleted(ctx context.Context, id int64, resourceID int64) error
}
