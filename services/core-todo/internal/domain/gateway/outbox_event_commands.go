//go:generate mockgen -destination=mock/outbox_event_commands_mock.go -source=outbox_event_commands.go -package mock

package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
)

type OutboxEventCommandsGateway interface {
	Create(ctx context.Context, in *input.CreateOutboxEvent) error
	MarkPublished(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64, errMsg string) error
}
