package output

import "time"

type OutboxEventStatus string

const (
	OutboxEventStatusPending OutboxEventStatus = "PENDING"
	OutboxEventStatusSent    OutboxEventStatus = "SENT"
	OutboxEventStatusFailed  OutboxEventStatus = "FAILED"
)

type OutboxEvent struct {
	ID          int64
	EventName   string
	RoutingKey  string
	Payload     []byte
	Status      OutboxEventStatus
	RetryCount  int
	LastError   *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt *time.Time
}
