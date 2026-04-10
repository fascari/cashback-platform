package redislock

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

var (
	// ErrLockNotAcquired is returned by Acquire when the lock is already held by another caller.
	ErrLockNotAcquired = errors.New("lock already held by another process")

	//go:embed acquire_lock.lua
	acquireLockScript string

	acquireScript = goredis.NewScript(acquireLockScript)
)

// Acquire obtains a distributed lock identified by lockKey using Redis SET NX
// and an incrementing fence token stored at fenceKey. ttlMs is the lock TTL in
// milliseconds. On success it returns the fence token and a release function the
// caller must invoke when the critical section ends. Returns ErrLockNotAcquired
// if another process already holds the lock.
func Acquire(ctx context.Context, client *goredis.Client, lockKey, fenceKey, lockValue string, ttlMs int) (fenceToken int64, release func(), err error) {
	token, err := acquireScript.Run(ctx, client, []string{lockKey, fenceKey}, lockValue, ttlMs).Int64()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, nil, ErrLockNotAcquired
		}
		return 0, nil, fmt.Errorf("acquire lock %s: %w", lockKey, err)
	}
	return token, func() { client.Del(ctx, lockKey) }, nil
}
