package input

import "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"

type CreateAuditLog struct {
	EventName  string
	EntityType string
	EntityID   int64
	ActorID    entity.UserID
	Payload    []byte
}
