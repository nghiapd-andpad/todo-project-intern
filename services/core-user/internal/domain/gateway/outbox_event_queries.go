//go:generate mockgen -destination=mock/outbox_event_queries_mock.go -source=outbox_event_queries.go -package mock

package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway/output"
)

type OutboxEventQueriesGateway interface {
	// FindClaimable selects PENDING/FAILED events, and PROCESSING events stuck past StuckThreshold, using SELECT FOR UPDATE SKIP LOCKED.
	FindClaimable(ctx context.Context, in *input.FindClaimableOutboxEvents) ([]*output.OutboxEvent, error)
}
