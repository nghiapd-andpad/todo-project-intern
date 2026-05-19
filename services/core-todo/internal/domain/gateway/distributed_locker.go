package gateway

import (
	"context"
	"time"
)

type DistributedLocker interface {
	// TryLock attempts to acquire a distributed lock for the given key.
	// The lock expires automatically after ttl if not released.
	//
	// Returns:
	//   - (true, release, nil) if the lock was acquired.
	//   - (false, nil, nil) if the lock is already held.
	//   - (false, nil, err) if an infrastructure error occurred.
	//
	// The caller must invoke release after the protected work completes.
	TryLock(ctx context.Context, key string, ttl time.Duration) (acquired bool, release func(), err error)
}
