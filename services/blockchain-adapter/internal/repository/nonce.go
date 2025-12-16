package repository

import (
	"context"

	"github.com/cashback-platform/services/blockchain-adapter/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	NonceRepository interface {
		GetAndIncrement(ctx context.Context, walletAddress string) (int64, error)
		GetCurrentNonce(ctx context.Context, walletAddress string) (int64, error)
	}

	nonceRepository struct {
		db *gorm.DB
	}
)

func NewNonceRepository(db *gorm.DB) NonceRepository {
	return &nonceRepository{db: db}
}

func (r *nonceRepository) GetAndIncrement(ctx context.Context, walletAddress string) (int64, error) {
	var nonce domain.WalletNonce

	// Use upsert with locking to ensure atomic nonce increment
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Try to find existing record with lock
		result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("wallet_address = ?", walletAddress).
			First(&nonce)

		if result.Error == gorm.ErrRecordNotFound {
			// Create new record
			nonce = domain.WalletNonce{
				ID:            uuid.New(),
				WalletAddress: walletAddress,
				CurrentNonce:  0,
			}
			if err := tx.Create(&nonce).Error; err != nil {
				return err
			}
		} else if result.Error != nil {
			return result.Error
		}

		// Increment nonce
		currentNonce := nonce.CurrentNonce
		nonce.CurrentNonce++

		if err := tx.Save(&nonce).Error; err != nil {
			return err
		}

		nonce.CurrentNonce = currentNonce // Return the nonce before increment
		return nil
	})

	if err != nil {
		return 0, err
	}

	return nonce.CurrentNonce, nil
}

func (r *nonceRepository) GetCurrentNonce(ctx context.Context, walletAddress string) (int64, error) {
	var nonce domain.WalletNonce
	if err := r.db.WithContext(ctx).Where("wallet_address = ?", walletAddress).First(&nonce).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return nonce.CurrentNonce, nil
}
