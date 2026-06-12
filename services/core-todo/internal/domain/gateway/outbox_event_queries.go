//go:generate mockgen -destination=mock/outbox_event_queries_mock.go -source=outbox_event_queries.go -package mock

package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/output"
)

type OutboxEventQueriesGateway interface {
	FindPending(ctx context.Context, in *input.ListPendingOutboxEvents) ([]*output.OutboxEvent, error)
}
