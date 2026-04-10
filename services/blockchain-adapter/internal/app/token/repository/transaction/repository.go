package transaction

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/cashback-platform/services/blockchain-adapter/internal/app/token/domain"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) Create(ctx context.Context, tx domain.BlockchainTransaction) (domain.BlockchainTransaction, error) {
	m := new(fromDomain(tx))
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return domain.BlockchainTransaction{}, fmt.Errorf("create transaction: %w", err)
	}
	return m.toDomain(), nil
}

func (r Repository) FindByIdempotencyKey(ctx context.Context, key uuid.UUID) (*domain.BlockchainTransaction, error) {
	m := new(transactionModel)
	if err := r.db.WithContext(ctx).Where("idempotency_key = ?", key).First(m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("find transaction by idempotency key: %w", err)
	}
	return new(m.toDomain()), nil
}

func (r Repository) UpdateStatus(ctx context.Context, id int64, status domain.TransactionStatus, txHash string, blockNumber int64) error {
	updates := map[string]any{
		"status":           status,
		"transaction_hash": txHash,
		"block_number":     blockNumber,
		"updated_at":       time.Now().UTC(),
	}
	if err := r.db.WithContext(ctx).Model(&transactionModel{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("update transaction status: %w", err)
	}
	return nil
}

func (r Repository) MarkFailed(ctx context.Context, id int64, errCode, errMsg string) error {
	updates := map[string]any{
		"status":        domain.TransactionStatusFailed,
		"error_code":    errCode,
		"error_message": errMsg,
		"updated_at":    time.Now().UTC(),
	}
	if err := r.db.WithContext(ctx).Model(&transactionModel{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("mark transaction failed: %w", err)
	}
	return nil
}
