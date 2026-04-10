package nonce

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/cashback-platform/kit/redislock"
	redisclient "github.com/cashback-platform/services/blockchain-adapter/internal/infra/redis"
)

const lockTTLMs = 30_000

var ErrStaleLockToken = errors.New("stale fencing token: lock was acquired by a newer holder")

type Repository struct {
	db    *gorm.DB
	redis *redisclient.Client
}

func NewRepository(db *gorm.DB, redis *redisclient.Client) Repository {
	return Repository{db: db, redis: redis}
}

func (r Repository) Increment(ctx context.Context, walletAddress string) (int64, error) {
	lockKey := "nonce:lock:" + walletAddress
	fenceKey := "nonce:fence:" + walletAddress
	lockValue := strconv.FormatInt(time.Now().UnixNano(), 10)

	fenceToken, release, err := redislock.Acquire(ctx, r.redis.Inner(), lockKey, fenceKey, lockValue, lockTTLMs)
	if err != nil {
		return 0, err
	}
	defer release()

	nonce, err := r.incrementInTx(ctx, walletAddress, fenceToken)
	if err != nil {
		return 0, fmt.Errorf("increment nonce for %s: %w", walletAddress, err)
	}
	return nonce, nil
}

func (r Repository) incrementInTx(ctx context.Context, walletAddress string, fenceToken int64) (int64, error) {
	var current int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		m := new(walletNonceModel)
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("wallet_address = ?", walletAddress).
			First(m).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return tx.Create(&walletNonceModel{
				WalletAddress: walletAddress,
				CurrentNonce:  1,
				FenceToken:    fenceToken,
			}).Error
		}

		if fenceToken <= m.FenceToken {
			return ErrStaleLockToken
		}

		current = m.CurrentNonce
		m.CurrentNonce++
		m.FenceToken = fenceToken
		return tx.Save(m).Error
	})
	return current, err
}

func (r Repository) CurrentNonce(ctx context.Context, walletAddress string) (int64, error) {
	m := new(walletNonceModel)
	if err := r.db.WithContext(ctx).Where("wallet_address = ?", walletAddress).First(m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, fmt.Errorf("get nonce for %s: %w", walletAddress, err)
	}
	return m.CurrentNonce, nil
}

func (r Repository) SyncFromChain(ctx context.Context, walletAddress string, nonce int64) error {
	m := walletNonceModel{
		WalletAddress: walletAddress,
		CurrentNonce:  nonce,
	}
	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "wallet_address"}},
			DoUpdates: clause.AssignmentColumns([]string{"current_nonce"}),
		}).
		Create(&m).Error; err != nil {
		return fmt.Errorf("sync nonce from chain for %s: %w", walletAddress, err)
	}
	return nil
}
