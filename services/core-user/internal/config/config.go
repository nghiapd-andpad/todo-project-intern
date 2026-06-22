package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBHost     string `envconfig:"DB_HOST" default:"localhost"`
	DBPort     string `envconfig:"DB_PORT" default:"3306"`
	DBUser     string `envconfig:"DB_USER" default:"root"`
	DBPassword string `envconfig:"DB_PASSWORD" default:"root"`
	DBName     string `envconfig:"DB_NAME" default:"user_db"`
	ServerPort string `envconfig:"SERVER_PORT" default:"50052"`
	AppEnv     string `envconfig:"APP_ENV" default:"development"`

	JWTSecret      string `envconfig:"JWT_SECRET" default:"secret_key"`
	JWTExpireHours int    `envconfig:"JWT_EXPIRATION_HRS" default:"24"`

	LogLevel  string `envconfig:"LOG_LEVEL" default:"debug"`
	LogFormat string `envconfig:"LOG_FORMAT" default:"console"`

	// Redis
	RedisHost     string `envconfig:"REDIS_HOST" default:"localhost"`
	RedisPort     string `envconfig:"REDIS_PORT" default:"6379"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`
	RedisDB       int    `envconfig:"REDIS_DB" default:"0"`

	// RabbitMQ — core-user connects to the same broker as core-todo
	// to consume events from the todo.notification.queue
	RabbitMQHost     string `envconfig:"RABBITMQ_HOST" default:"localhost"`
	RabbitMQPort     string `envconfig:"RABBITMQ_PORT" default:"5672"`
	RabbitMQUser     string `envconfig:"RABBITMQ_USER" default:"root"`
	RabbitMQPassword string `envconfig:"RABBITMQ_PASSWORD" default:"root"`

	// core-todo publishes to this exchange, core-user subscribes
	RabbitMQTodoExchange              string `envconfig:"RABBITMQ_TODO_EXCHANGE" default:"todo.events"`
	RabbitMQNotificationQueue         string `envconfig:"RABBITMQ_NOTIFICATION_QUEUE" default:"todo.notification.queue"`
	RabbitMQNotificationRoutingKey    string `envconfig:"RABBITMQ_NOTIFICATION_ROUTING_KEY" default:"todo.assigned"`
	RabbitMQNotificationPrefetchCount int    `envconfig:"RABBITMQ_NOTIFICATION_PREFETCH_COUNT" default:"10"`

	// Outbox publisher worker
	SchedulerEnabled bool `envconfig:"SCHEDULER_ENABLED" default:"false"`

	OutboxPublisherCron      string `envconfig:"OUTBOX_PUBLISHER_CRON" default:"*/1 * * * *"`
	OutboxPublisherBatchSize int    `envconfig:"OUTBOX_PUBLISHER_BATCH_SIZE" default:"100"`
	OutboxPublisherMaxRetry  int    `envconfig:"OUTBOX_PUBLISHER_MAX_RETRY" default:"5"`

	// core-user publishes here
	RabbitMQUserExchange string `envconfig:"RABBITMQ_USER_EXCHANGE" default:"user.events"`
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
