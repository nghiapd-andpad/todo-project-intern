// Package redis provides Redis client initialization and Redis-based infrastructure implementations.
package redis

import (
	"context"
	"fmt"
	"net"

	redisv9 "github.com/redis/go-redis/v9"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
)

func NewClient(cfg *config.Config) (*redisv9.Client, func(), error) {
	if cfg == nil {
		return nil, nil, fmt.Errorf("config is nil")
	}

	addr := net.JoinHostPort(cfg.RedisHost, cfg.RedisPort)

	client := redisv9.NewClient(&redisv9.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	cleanup := func() {
		_ = client.Close()
	}

	return client, cleanup, nil
}
