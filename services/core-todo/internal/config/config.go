// Package config provides configuration loading and management for the application.
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

	// Log
	LogLevel  string `envconfig:"LOG_LEVEL" default:"debug"`
	LogFormat string `envconfig:"LOG_FORMAT" default:"console"`

	// Scheduler
	SchedulerEnabled            bool          `envconfig:"SCHEDULER_ENABLED" default:"false"`
	TodoOverdueMarkerCron       string        `envconfig:"TODO_OVERDUE_MARKER_CRON" default:"*/5 * * * *"`
	TodoOverdueMarkerBatchSize  int           `envconfig:"TODO_OVERDUE_MARKER_BATCH_SIZE" default:"500"`
	TodoOverdueMarkerMaxBatches int           `envconfig:"TODO_OVERDUE_MARKER_MAX_BATCHES" default:"20"`
	TodoOverdueMarkerLockKey    string        `envconfig:"TODO_OVERDUE_MARKER_LOCK_KEY" default:"job:mark-overdue-todos"`
	TodoOverdueMarkerLockTTL    time.Duration `envconfig:"TODO_OVERDUE_MARKER_LOCK_TTL" default:"10m"`
	TodoOverdueMarkerBatchSleep time.Duration `envconfig:"TODO_OVERDUE_MARKER_BATCH_SLEEP" default:"100ms"`

	// Redis
	RedisHost     string `envconfig:"REDIS_HOST" default:"localhost"`
	RedisPort     string `envconfig:"REDIS_PORT" default:"6379"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`
	RedisDB       int    `envconfig:"REDIS_DB" default:"0"`

	// Feature flags
	TodoBlacklistEnabled bool     `envconfig:"TODO_BLACKLIST_ENABLED" default:"false"`
	TodoTitleBlacklist   []string `envconfig:"TODO_TITLE_BLACKLIST" default:"spam,troll"`
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
