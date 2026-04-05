package mintrequest

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/cashback-platform/kit/gormtx"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/domain"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) Create(ctx context.Context, request domain.MintRequest) (domain.MintRequest, error) {
	m := new(fromDomain(request))
	db := r.db.WithContext(ctx)
	if tx := gormtx.ExtractTx(ctx); tx != nil {
		db = tx.WithContext(ctx)
	}
	if err := db.Create(m).Error; err != nil {
		return domain.MintRequest{}, err
	}
	return m.toDomain(), nil
}

func (r Repository) FindByID(ctx context.Context, id int64) (domain.MintRequest, error) {
	m := new(mintRequestModel)
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(m).Error; err != nil {
		return domain.MintRequest{}, err
	}
	return m.toDomain(), nil
}

func (r Repository) FindByCashbackID(ctx context.Context, cashbackID int64) (domain.MintRequest, error) {
	m := new(mintRequestModel)
	if err := r.db.WithContext(ctx).Where("cashback_id = ?", cashbackID).First(m).Error; err != nil {
		return domain.MintRequest{}, err
	}
	return m.toDomain(), nil
}

func (r Repository) FindByIdempotencyKey(ctx context.Context, key uuid.UUID) (domain.MintRequest, error) {
	m := new(mintRequestModel)
	if err := r.db.WithContext(ctx).Where("idempotency_key = ?", key).First(m).Error; err != nil {
		return domain.MintRequest{}, err
	}
	return m.toDomain(), nil
}

func (r Repository) Update(ctx context.Context, request domain.MintRequest) error {
	m := fromDomain(request)
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r Repository) UpdateStatus(ctx context.Context, id int64, status domain.MintRequestStatus) error {
	return r.db.WithContext(ctx).Model(&mintRequestModel{}).Where("id = ?", id).Update("status", status).Error
}

func (r Repository) FindFailedRetryable(ctx context.Context, limit int) ([]domain.MintRequest, error) {
	var models []mintRequestModel
	now := time.Now().UTC()
	err := r.db.WithContext(ctx).
		Where("status = ? AND next_retry_at <= ? AND retry_count < max_retries", domain.MintRequestStatusFailed, now).
		Order("next_retry_at ASC").
		Limit(limit).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	requests := make([]domain.MintRequest, len(models))
	for i, m := range models {
		requests[i] = m.toDomain()
	}
	return requests, nil
}

func (r Repository) MarkCompleted(ctx context.Context, id int64, txHash string, blockNumber int64) error {
	completedAt := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&mintRequestModel{}).Where("id = ?", id).Updates(map[string]any{
		"status":           domain.MintRequestStatusCompleted,
		"transaction_hash": txHash,
		"block_number":     blockNumber,
		"completed_at":     &completedAt,
	}).Error
}

func (r Repository) MarkFailed(ctx context.Context, id int64, errorCode, errorMessage string, nextRetryAt *time.Time) error {
	return r.db.WithContext(ctx).Model(&mintRequestModel{}).Where("id = ?", id).Updates(map[string]any{
		"status":        domain.MintRequestStatusFailed,
		"error_code":    errorCode,
		"error_message": errorMessage,
		"next_retry_at": nextRetryAt,
		"retry_count":   gorm.Expr("retry_count + 1"),
	}).Error
}
