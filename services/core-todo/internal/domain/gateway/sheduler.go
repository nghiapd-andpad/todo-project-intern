package gateway

import (
	"context"
	"time"
)

type Scheduler interface {
	ScheduleOnce(ctx context.Context, jobTag string, duration time.Duration, function any, parameters ...any) error
}
