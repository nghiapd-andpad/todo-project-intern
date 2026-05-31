package mapper

import (
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/output"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/model"
)

func IdempotencyRecordFromModel(m *model.IdempotencyKey) *output.IdempotencyRecord {
	if m == nil {
		return nil
	}

	return &output.IdempotencyRecord{
		ID:             m.ID,
		UserID:         entity.UserID(m.UserID),
		Operation:      m.Operation,
		IdempotencyKey: m.IdempotencyKey,
		Status:         output.IdempotencyStatus(m.Status),
		ResourceID:     m.ResourceID,
		ExpiresAt:      m.ExpiresAt,
	}
}
