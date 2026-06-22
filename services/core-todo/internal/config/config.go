package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBHost     string `envconfig:"DB_HOST" default:"localhost"`
	DBPort     string `envconfig:"DB_PORT" default:"3306"`
	DBUser     string `envconfig:"DB_USER" default:"root"`
	DBPassword string `envconfig:"DB_PASSWORD" default:"root"`
	DBName     string `envconfig:"DB_NAME" default:"todos"`
	ServerPort string `envconfig:"SERVER_PORT" default:"50051"`
	AppEnv     string `envconfig:"APP_ENV" default:"development"`

	LogLevel  string `envconfig:"LOG_LEVEL" default:"debug"`
	LogFormat string `envconfig:"LOG_FORMAT" default:"console"`

	SchedulerEnabled bool `envconfig:"SCHEDULER_ENABLED" default:"false"`

	// Mark overdue todo job
	TodoOverdueMarkerCron       string        `envconfig:"TODO_OVERDUE_MARKER_CRON" default:"*/5 * * * *"`
	TodoOverdueMarkerBatchSize  int           `envconfig:"TODO_OVERDUE_MARKER_BATCH_SIZE" default:"500"`
	TodoOverdueMarkerMaxBatches int           `envconfig:"TODO_OVERDUE_MARKER_MAX_BATCHES" default:"20"`
	TodoOverdueMarkerLockKey    string        `envconfig:"TODO_OVERDUE_MARKER_LOCK_KEY" default:"job:mark-overdue-todos"`
	TodoOverdueMarkerLockTTL    time.Duration `envconfig:"TODO_OVERDUE_MARKER_LOCK_TTL" default:"10m"`
	TodoOverdueMarkerBatchSleep time.Duration `envconfig:"TODO_OVERDUE_MARKER_BATCH_SLEEP" default:"100ms"`

	// Soft deleted cleanup job
	TodoSoftDeletedCleanupCron          string        `envconfig:"TODO_SOFT_DELETED_CLEANUP_CRON" default:"0 0 * * *"`
	TodoSoftDeletedCleanupBatchSize     int           `envconfig:"TODO_SOFT_DELETED_CLEANUP_BATCH_SIZE" default:"500"`
	TodoSoftDeletedCleanupMaxBatches    int           `envconfig:"TODO_SOFT_DELETED_CLEANUP_MAX_BATCHES" default:"20"`
	TodoSoftDeletedCleanupRetentionDays int           `envconfig:"TODO_SOFT_DELETED_CLEANUP_RETENTION_DAYS" default:"30"`
	TodoSoftDeletedCleanupLockKey       string        `envconfig:"TODO_SOFT_DELETED_CLEANUP_LOCK_KEY" default:"job:cleanup-soft-deleted-todos"`
	TodoSoftDeletedCleanupLockTTL       time.Duration `envconfig:"TODO_SOFT_DELETED_CLEANUP_LOCK_TTL" default:"10m"`
	TodoSoftDeletedCleanupBatchSleep    time.Duration `envconfig:"TODO_SOFT_DELETED_CLEANUP_BATCH_SLEEP" default:"100ms"`

	// Outbox publisher worker (no Redis lock — uses SKIP LOCKED instead)
	OutboxPublisherCron      string `envconfig:"OUTBOX_PUBLISHER_CRON" default:"*/1 * * * *"`
	OutboxPublisherBatchSize int    `envconfig:"OUTBOX_PUBLISHER_BATCH_SIZE" default:"100"`
	OutboxPublisherMaxRetry  int    `envconfig:"OUTBOX_PUBLISHER_MAX_RETRY" default:"5"`

	// Redis
	RedisHost     string `envconfig:"REDIS_HOST" default:"localhost"`
	RedisPort     string `envconfig:"REDIS_PORT" default:"6379"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`
	RedisDB       int    `envconfig:"REDIS_DB" default:"0"`

	// Feature flags
	TodoBlacklistEnabled bool     `envconfig:"TODO_BLACKLIST_ENABLED" default:"false"`
	TodoTitleBlacklist   []string `envconfig:"TODO_TITLE_BLACKLIST" default:"spam,troll"`

	// RabbitMQ
	RabbitMQHost     string `envconfig:"RABBITMQ_HOST" default:"localhost"`
	RabbitMQPort     string `envconfig:"RABBITMQ_PORT" default:"5672"`
	RabbitMQUser     string `envconfig:"RABBITMQ_USER" default:"root"`
	RabbitMQPassword string `envconfig:"RABBITMQ_PASSWORD" default:"root"`

	// core-todo publishes to this exchange
	RabbitMQTodoExchange string `envconfig:"RABBITMQ_TODO_EXCHANGE" default:"todo.events"`

	// core-user publishes to this exchange, core-todo subscribes
	RabbitMQUserExchange            string `envconfig:"RABBITMQ_USER_EXCHANGE" default:"user.events"`
	RabbitMQUserEventsQueue         string `envconfig:"RABBITMQ_USER_EVENTS_QUEUE" default:"todo.user-events.queue"`
	RabbitMQUserEventsRoutingKey    string `envconfig:"RABBITMQ_USER_EVENTS_ROUTING_KEY" default:"user.*"`
	RabbitMQUserEventsPrefetchCount int    `envconfig:"RABBITMQ_USER_EVENTS_PREFETCH_COUNT" default:"10"`

	RabbitMQNotificationQueue      string `envconfig:"RABBITMQ_NOTIFICATION_QUEUE" default:"todo.notification.queue"`
	RabbitMQNotificationRoutingKey string `envconfig:"RABBITMQ_NOTIFICATION_ROUTING_KEY" default:"todo.assigned"`
}

func New() (*Config, error) {
	env := os.Getenv("APP_ENV")

	var envFile string
	switch env {
	case "production":
		envFile = ".env.production"
	default:
		envFile = ".env"
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("no %s file found, fallback to system env", envFile)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("load config error: %w", err)
	}

	return &cfg, nil
}
