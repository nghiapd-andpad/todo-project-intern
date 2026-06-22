package input

import "time"

type CreateOutboxEvent struct {
	EventName  string
	RoutingKey string
	Payload    []byte
}

type FindClaimableOutboxEvents struct {
	BatchSize      int
	MaxRetry       int
	StuckThreshold time.Duration
}
