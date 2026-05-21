package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	redisv9 "github.com/redis/go-redis/v9"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/gateway"
)

const releaseLockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else
	return 0
end
`

type DistributedLocker struct {
	client *redisv9.Client
}

var _ gateway.DistributedLocker = (*DistributedLocker)(nil)

func NewDistributedLocker(client *redisv9.Client) *DistributedLocker {
	return &DistributedLocker{
		client: client,
	}
}

func (l *DistributedLocker) TryLock(ctx context.Context, key string, ttl time.Duration) (bool, func(), error) {
	if l.client == nil {
		return false, nil, fmt.Errorf("redis client is nil")
	}
	if key == "" {
		return false, nil, fmt.Errorf("lock key is empty")
	}
	if ttl <= 0 {
		return false, nil, fmt.Errorf("lock ttl must be greater than zero")
	}

	token := uuid.NewString()

	acquired, err := l.client.SetNX(ctx, key, token, ttl).Result()
	if err != nil {
		return false, nil, fmt.Errorf("redis acquire lock: %w", err)
	}

	if !acquired {
		return false, nil, nil
	}

	released := false

	release := func() {
		if released {
			return
		}
		released = true

		_ = l.client.Eval(
			context.Background(),
			releaseLockScript,
			[]string{key},
			token,
		).Err()
	}

	return true, release, nil
}
