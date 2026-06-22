package output

import "time"

type OutboxEventStatus string

const (
	OutboxEventStatusPending    OutboxEventStatus = "PENDING"
	OutboxEventStatusProcessing OutboxEventStatus = "PROCESSING"
	OutboxEventStatusPublished  OutboxEventStatus = "PUBLISHED"
	OutboxEventStatusFailed     OutboxEventStatus = "FAILED"
	OutboxEventStatusDead       OutboxEventStatus = "DEAD"
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
