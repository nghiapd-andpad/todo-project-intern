package gateway

import "context"

type ProcessedEventGateway interface {
	// TryRecord attempts to log a processed event.
	// Returns true if the event was successfully recorded (new event), and false if the event was already processed (duplicate).
	TryRecord(ctx context.Context, hash string, consumerKey string) (bool, error)
}
