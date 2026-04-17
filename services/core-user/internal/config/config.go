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

	// Server Config
	ServerPort string `envconfig:"SERVER_PORT" default:"50052"`
	AppEnv     string `envconfig:"APP_ENV" default:"development"`

	// JWT Security Config (Bổ sung cho Auth)
	JWTSecret        string `envconfig:"JWT_SECRET" default:"secret_key"`
	JWTExpirationHrs int    `envconfig:"JWT_EXPIRATION_HRS" default:"24"`
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
