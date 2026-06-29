//go:generate mockgen -destination=mock/audit_log_commands_mock.go -source=audit_log_commands.go -package mock

package gateway

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
)

type AuditLogCommandsGateway interface {
	Create(ctx context.Context, in *input.CreateAuditLog) error
}
