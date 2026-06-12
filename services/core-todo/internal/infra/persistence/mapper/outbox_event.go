package mapper

import (
	gatewayoutput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/output"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/infra/persistence/model"
)

func OutboxEventToOutput(m *model.OutboxEvent) *gatewayoutput.OutboxEvent {
	if m == nil {
		return nil
	}

	return &gatewayoutput.OutboxEvent{
		ID:          m.ID,
		EventName:   m.EventName,
		RoutingKey:  m.RoutingKey,
		Payload:     []byte(m.Payload),
		Status:      gatewayoutput.OutboxEventStatus(m.Status),
		RetryCount:  m.RetryCount,
		LastError:   m.LastError,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		PublishedAt: m.PublishedAt,
	}
}

func OutboxEventsToOutput(models []*model.OutboxEvent) []*gatewayoutput.OutboxEvent {
	result := make([]*gatewayoutput.OutboxEvent, len(models))
	for i := range models {
		result[i] = OutboxEventToOutput(models[i])
	}
	return result
}
