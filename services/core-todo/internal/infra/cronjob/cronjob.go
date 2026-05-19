package cronjob

import (
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

const defaultLimitConcurrentJobs = 10

func NewGoCronScheduler(zapLogger *zap.Logger) (gocron.Scheduler, func(), error) {
	l := newLogger(zapLogger)

	scheduler, err := gocron.NewScheduler(
		gocron.WithLogger(l),
		gocron.WithLimitConcurrentJobs(defaultLimitConcurrentJobs, gocron.LimitModeWait),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create gocron scheduler: %w", err)
	}

	cleanup := func() {
		if err := scheduler.Shutdown(); err != nil {
			l.Error("failed to shutdown scheduler: %v", err)
		}
	}

	return scheduler, cleanup, nil
}
