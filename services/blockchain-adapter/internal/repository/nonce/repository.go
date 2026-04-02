package nonce

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/cashback-platform/services/blockchain-adapter/internal/domain"
	redisclient "github.com/cashback-platform/services/blockchain-adapter/internal/infra/redis"
)

const (
	lockTTLMs = 10_000

	// acquireLockScript uses a Lua transaction to guarantee that the lock acquisition and
	// fence token increment are atomic on the Redis side — no other caller can interleave.
	acquireLockScript = `
		local acquired = redis.call('SET', KEYS[1], ARGV[1], 'NX', 'PX', ARGV[2])
		if acquired then
			return redis.call('INCR', KEYS[2])
		end
		return nil
		`
)

var (
	ErrLockNotAcquired = errors.New("nonce lock already held by another process")
	ErrStaleLockToken  = errors.New("stale fencing token: lock was acquired by a newer holder")

	acquireScript = goredis.NewScript(acquireLockScript)
)

type Repository struct {
	db    *gorm.DB
	redis *redisclient.Client
}

func NewRepository(db *gorm.DB, redis *redisclient.Client) *Repository {
	return &Repository{db: db, redis: redis}
}

func (r *Repository) Increment(ctx context.Context, walletAddress string) (int64, error) {
	fenceToken, err := r.acquireLock(ctx, walletAddress)
	if err != nil {
		return 0, err
	}
	defer r.redis.Inner().Del(ctx, "nonce:lock:"+walletAddress)

	nonce, err := r.incrementInTx(ctx, walletAddress, fenceToken)
	if err != nil {
		return 0, fmt.Errorf("increment nonce for %s: %w", walletAddress, err)
	}
	return nonce, nil
}

func (r *Repository) acquireLock(ctx context.Context, walletAddress string) (int64, error) {
	lockKey := "nonce:lock:" + walletAddress
	fenceKey := "nonce:fence:" + walletAddress
	lockValue := strconv.FormatInt(time.Now().UnixNano(), 10)

	fenceToken, err := acquireScript.Run(ctx, r.redis.Inner(), []string{lockKey, fenceKey}, lockValue, lockTTLMs).Int64()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return 0, ErrLockNotAcquired
		}
		return 0, fmt.Errorf("acquire nonce lock: %w", err)
	}
	return fenceToken, nil
}

func (r *Repository) incrementInTx(ctx context.Context, walletAddress string, fenceToken int64) (int64, error) {
	var current int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var nonce domain.WalletNonce
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("wallet_address = ?", walletAddress).
			First(&nonce).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return tx.Create(&domain.WalletNonce{
				WalletAddress: walletAddress,
				CurrentNonce:  1,
				FenceToken:    fenceToken,
			}).Error
		}

		if fenceToken <= nonce.FenceToken {
			return ErrStaleLockToken
		}

		current = nonce.CurrentNonce
		nonce.CurrentNonce++
		nonce.FenceToken = fenceToken
		return tx.Save(&nonce).Error
	})
	return current, err
}

func (r *Repository) CurrentNonce(ctx context.Context, walletAddress string) (int64, error) {
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
	if err := r.db.WithContext(ctx).
		Model(&domain.WalletNonce{}).
		Where("wallet_address = ?", walletAddress).
		Updates(map[string]any{"current_nonce": nonce}).Error; err != nil {
		return fmt.Errorf("sync nonce from chain for %s: %w", walletAddress, err)
	}
	return nil
}
