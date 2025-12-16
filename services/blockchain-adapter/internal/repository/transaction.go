package repository

import (
	"context"
	"time"

	"github.com/cashback-platform/services/blockchain-adapter/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	TransactionRepository interface {
		Create(ctx context.Context, tx *domain.BlockchainTransaction) error
		GetByID(ctx context.Context, id uuid.UUID) (*domain.BlockchainTransaction, error)
		GetByIdempotencyKey(ctx context.Context, key uuid.UUID) (*domain.BlockchainTransaction, error)
		GetByTransactionHash(ctx context.Context, hash string) (*domain.BlockchainTransaction, error)
		Update(ctx context.Context, tx *domain.BlockchainTransaction) error
		UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TransactionStatus) error
		MarkConfirmed(ctx context.Context, id uuid.UUID, blockNumber int64, gasUsed int64) error
		MarkFailed(ctx context.Context, id uuid.UUID, errorCode, errorMessage string) error
	}

	transactionRepository struct {
		db *gorm.DB
	}
)

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.BlockchainTransaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *transactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.BlockchainTransaction, error) {
	var tx domain.BlockchainTransaction
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&tx).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) GetByIdempotencyKey(ctx context.Context, key uuid.UUID) (*domain.BlockchainTransaction, error) {
	var tx domain.BlockchainTransaction
	if err := r.db.WithContext(ctx).Where("idempotency_key = ?", key).First(&tx).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) GetByTransactionHash(ctx context.Context, hash string) (*domain.BlockchainTransaction, error) {
	var tx domain.BlockchainTransaction
	if err := r.db.WithContext(ctx).Where("transaction_hash = ?", hash).First(&tx).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) Update(ctx context.Context, tx *domain.BlockchainTransaction) error {
	return r.db.WithContext(ctx).Save(tx).Error
}

func (r *transactionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TransactionStatus) error {
	return r.db.WithContext(ctx).Model(&domain.BlockchainTransaction{}).Where("id = ?", id).Update("status", status).Error
}

func (r *transactionRepository) MarkConfirmed(ctx context.Context, id uuid.UUID, blockNumber int64, gasUsed int64) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&domain.BlockchainTransaction{}).Where("id = ?", id).Updates(map[string]any{
		"status":       domain.TransactionStatusConfirmed,
		"block_number": blockNumber,
		"gas_used":     gasUsed,
		"confirmed_at": &now,
	}).Error
}

func (r *transactionRepository) MarkFailed(ctx context.Context, id uuid.UUID, errorCode, errorMessage string) error {
	return r.db.WithContext(ctx).Model(&domain.BlockchainTransaction{}).Where("id = ?", id).Updates(map[string]any{
		"status":        domain.TransactionStatusFailed,
		"error_code":    errorCode,
		"error_message": errorMessage,
	}).Error
}
