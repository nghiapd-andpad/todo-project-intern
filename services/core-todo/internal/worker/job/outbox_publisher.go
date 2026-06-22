package job

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
	gatewayinput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/input"
	gatewayoutput "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway/output"
	logutil "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/utils/logger"
)

const outboxPublisherJobName = "outbox_publisher"

// OutboxPublisherJob claims a batch of PENDING/FAILED outbox events using SELECT FOR UPDATE SKIP LOCKED, publishes each to RabbitMQ, then marks them PUBLISHED or FAILED/DEAD.
// Multiple instances can run concurrently — SKIP LOCKED ensures each instance processes a distinct batch with no duplicates and no Redis needed.
type OutboxPublisherJob struct {
	transactor     gateway.Transactor
	outboxQueries  gateway.OutboxEventQueriesGateway
	outboxCommands gateway.OutboxEventCommandsGateway
	publisher      gateway.EventPublisher
	cfg            *config.Config
	logger         *zap.Logger
}

func NewOutboxPublisherJob(
	transactor gateway.Transactor,
	outboxQueries gateway.OutboxEventQueriesGateway,
	outboxCommands gateway.OutboxEventCommandsGateway,
	publisher gateway.EventPublisher,
	cfg *config.Config,
	zapLogger *zap.Logger,
) *OutboxPublisherJob {
	if zapLogger == nil {
		zapLogger = zap.NewNop()
	}

	return &OutboxPublisherJob{
		transactor:     transactor,
		outboxQueries:  outboxQueries,
		outboxCommands: outboxCommands,
		publisher:      publisher,
		cfg:            cfg,
		logger: zapLogger.With(
			zap.String("component", "cronjob"),
			zap.String("job_name", outboxPublisherJobName),
		),
	}
}

func (j *OutboxPublisherJob) Name() string { return outboxPublisherJobName }

func (j *OutboxPublisherJob) Cron() string {
	if j.cfg == nil {
		return ""
	}
	return j.cfg.OutboxPublisherCron
}

func (j *OutboxPublisherJob) Run() error {
	start := time.Now()
	ctx := logutil.ToContext(context.Background(), j.logger)

	logutil.Info(ctx, "outbox publisher job started")

	// StuckThreshold > Worst Case Processing Time
	// worst case ≈ BatchSize × per-event-timeout(10s) × safety-factor(2)
	stuckThreshold := time.Duration(j.cfg.OutboxPublisherBatchSize) * 10 * time.Second * 2

	var events []*gatewayoutput.OutboxEvent

	err := j.transactor.Transaction(ctx, func(txCtx context.Context) error {
		claimed, err := j.outboxQueries.FindClaimable(txCtx, &gatewayinput.FindClaimableOutboxEvents{
			BatchSize:      j.cfg.OutboxPublisherBatchSize,
			MaxRetry:       j.cfg.OutboxPublisherMaxRetry,
			StuckThreshold: stuckThreshold,
		})
		if err != nil {
			return err
		}
		if len(claimed) == 0 {
			return nil
		}

		ids := make([]int64, len(claimed))
		for i, ev := range claimed {
			ids[i] = ev.ID
		}

		if err := j.outboxCommands.MarkProcessing(txCtx, ids); err != nil {
			return err
		}

		events = claimed
		return nil
	})
	if err != nil {
		return fmt.Errorf("claim outbox events: %w", err)
	}

	if len(events) == 0 {
		logutil.Info(ctx, "outbox publisher job: no events to process")
		return nil
	}

	var published, failed, dead int

	for _, ev := range events {
		if err := j.processOne(ctx, ev); err != nil {
			if ev.RetryCount >= j.cfg.OutboxPublisherMaxRetry {
				dead++
			} else {
				failed++
			}
			continue
		}
		published++
	}

	logutil.Info(ctx, "outbox publisher job completed",
		zap.Int("total", len(events)),
		zap.Int("published", published),
		zap.Int("failed", failed),
		zap.Int("dead", dead),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}

func (j *OutboxPublisherJob) processOne(ctx context.Context, ev *gatewayoutput.OutboxEvent) error {
	eventCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if ev.RetryCount > j.cfg.OutboxPublisherMaxRetry {
		errMsg := fmt.Sprintf("exceeded max retry count (%d)", j.cfg.OutboxPublisherMaxRetry)
		if dbErr := j.outboxCommands.MarkDead(ctx, ev.ID, errMsg); dbErr != nil {
			j.logger.Error("mark outbox event dead failed",
				zap.Int64("event_id", ev.ID), zap.Error(dbErr))
		}
		j.logger.Warn("outbox event marked dead",
			zap.Int64("event_id", ev.ID),
			zap.String("event_name", ev.EventName),
			zap.Int("retry_count", ev.RetryCount),
		)
		return fmt.Errorf("event %d is dead", ev.ID)
	}

	if err := j.publisher.Publish(eventCtx, ev.RoutingKey, ev.Payload); err != nil {
		j.logger.Error("publish outbox event failed",
			zap.Int64("event_id", ev.ID),
			zap.String("routing_key", ev.RoutingKey),
			zap.Int("retry_count", ev.RetryCount),
			zap.Error(err),
		)
		if dbErr := j.outboxCommands.MarkFailed(ctx, ev.ID, err.Error()); dbErr != nil {
			j.logger.Error("mark outbox event failed in db",
				zap.Int64("event_id", ev.ID), zap.Error(dbErr))
		}
		return err
	}

	if err := j.outboxCommands.MarkPublished(ctx, ev.ID); err != nil {
		var appErr *entity.AppError

		if errors.As(err, &appErr) &&
			appErr.Code == entity.ErrAlreadyHandled {

			j.logger.Warn(
				"outbox event already handled elsewhere",
				zap.Int64("event_id", ev.ID),
			)
			return nil
		}

		return err
	}

	j.logger.Debug("outbox event published",
		zap.Int64("event_id", ev.ID),
		zap.String("event_name", ev.EventName),
		zap.String("routing_key", ev.RoutingKey),
	)

	return nil
}
