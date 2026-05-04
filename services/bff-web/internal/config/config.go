// Package config loads configuration from environment variables and .env files.
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Server
	ServerPort string `envconfig:"SERVER_PORT" default:"8080"`
	AppEnv     string `envconfig:"APP_ENV" default:"development"`

	// gRPC endpoints
	TodoServiceAddr string `envconfig:"TODO_SERVICE_ADDR" default:"localhost:50051"`
	UserServiceAddr string `envconfig:"USER_SERVICE_ADDR" default:"localhost:50052"`

	// JWT secret
	JWTSecret string `envconfig:"JWT_SECRET" default:"secret_key"`

	// DataLoader
	DataLoaderWait      int `envconfig:"DATALOADER_WAIT_MS" default:"2"`      // ms
	DataLoaderBatchSize int `envconfig:"DATALOADER_BATCH_SIZE" default:"100"` // max batch size
}

func New() (*Config, error) {
	env := os.Getenv("APP_ENV")
	envFile := ".env"
	if env == "production" {
		envFile = ".env.production"
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("no %s file found, fallback to system env", envFile)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return &cfg, nil
}
