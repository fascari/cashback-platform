package nonce

import (
	"context"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/cashback-platform/services/blockchain-adapter/internal/domain"
	redisclient "github.com/cashback-platform/services/blockchain-adapter/internal/infra/redis"
)

var (
	// ErrLockNotAcquired is returned when the Redis lock is already held by another process.
	ErrLockNotAcquired = errors.New("nonce lock already held by another process")
	// ErrStaleLockToken is returned when the fencing token indicates a stale lock holder.
	ErrStaleLockToken = errors.New("stale fencing token: lock was acquired by a newer holder")
)

// acquireScript atomically acquires the lock and increments the fence token.
// Returns the new fence token, or nil if the lock is already held.
var acquireScript = goredis.NewScript(`
local acquired = redis.call('SET', KEYS[1], ARGV[1], 'NX', 'PX', ARGV[2])
if acquired then
    return redis.call('INCR', KEYS[2])
end
return nil
`)

const lockTTLMs = 10_000

type (
	NonceRepository interface {
		GetAndIncrement(ctx context.Context, walletAddress string) (int64, error)
		GetCurrentNonce(ctx context.Context, walletAddress string) (int64, error)
		SyncFromChain(ctx context.Context, walletAddress string, nonce int64) error
	}

	Repository struct {
		db    *gorm.DB
		redis *redisclient.Client
	}
)

// NewRepository creates a NonceRepository with Redis distributed locking.
func NewRepository(db *gorm.DB, redis *redisclient.Client) *Repository {
	return &Repository{db: db, redis: redis}
}

func (r *Repository) GetAndIncrement(ctx context.Context, walletAddress string) (int64, error) {
	lockKey := fmt.Sprintf("nonce:lock:%s", walletAddress)
	fenceKey := fmt.Sprintf("nonce:fence:%s", walletAddress)
	lockValue := fmt.Sprintf("%d", time.Now().UnixNano())

	result, err := acquireScript.Run(ctx, r.redis.Inner(), []string{lockKey, fenceKey}, lockValue, lockTTLMs).Int64()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, ErrLockNotAcquired
		}
		return 0, fmt.Errorf("acquire nonce lock: %w", err)
	}
	fenceToken := result

	defer r.redis.Inner().Del(ctx, lockKey)

	var nonce domain.WalletNonce
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("wallet_address = ?", walletAddress).
			First(&nonce)

		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			nonce = domain.WalletNonce{
				WalletAddress: walletAddress,
				CurrentNonce:  0,
				FenceToken:    fenceToken,
			}
			return tx.Create(&nonce).Error
		}
		if res.Error != nil {
			return res.Error
		}

		if fenceToken <= nonce.FenceToken {
			return ErrStaleLockToken
		}

		current := nonce.CurrentNonce
		nonce.CurrentNonce++
		nonce.FenceToken = fenceToken

		if err := tx.Save(&nonce).Error; err != nil {
			return err
		}
		nonce.CurrentNonce = current
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("increment nonce for %s: %w", walletAddress, err)
	}

	return nonce.CurrentNonce, nil
}

func (r *Repository) GetCurrentNonce(ctx context.Context, walletAddress string) (int64, error) {
	var nonce domain.WalletNonce
	if err := r.db.WithContext(ctx).Where("wallet_address = ?", walletAddress).First(&nonce).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, fmt.Errorf("get nonce for %s: %w", walletAddress, err)
	}
	return nonce.CurrentNonce, nil
}

func (r *Repository) SyncFromChain(ctx context.Context, walletAddress string, nonce int64) error {
	return r.db.WithContext(ctx).
		Model(&domain.WalletNonce{}).
		Where("wallet_address = ?", walletAddress).
		Updates(map[string]any{"current_nonce": nonce}).Error
}

