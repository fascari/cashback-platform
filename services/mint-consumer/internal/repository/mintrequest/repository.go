package mintrequest

import (
	"context"
	"time"

	"github.com/cashback-platform/services/mint-consumer/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) Create(ctx context.Context, request *domain.MintRequest) error {
	return r.db.WithContext(ctx).Create(request).Error
}

func (r Repository) FindByID(ctx context.Context, id uuid.UUID) (*domain.MintRequest, error) {
	request := new(domain.MintRequest)
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (r Repository) FindByCashbackID(ctx context.Context, cashbackID uuid.UUID) (*domain.MintRequest, error) {
	request := new(domain.MintRequest)
	if err := r.db.WithContext(ctx).Where("cashback_id = ?", cashbackID).First(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (r Repository) FindByIdempotencyKey(ctx context.Context, key uuid.UUID) (*domain.MintRequest, error) {
	request := new(domain.MintRequest)
	if err := r.db.WithContext(ctx).Where("idempotency_key = ?", key).First(request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (r Repository) Update(ctx context.Context, request *domain.MintRequest) error {
	return r.db.WithContext(ctx).Save(request).Error
}

func (r Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.MintRequestStatus) error {
	return r.db.WithContext(ctx).Model(&domain.MintRequest{}).Where("id = ?", id).Update("status", status).Error
}

func (r Repository) PendingRetries(ctx context.Context, limit int) ([]domain.MintRequest, error) {
	var requests []domain.MintRequest
	now := time.Now().UTC()
	err := r.db.WithContext(ctx).
		Where("status = ? AND next_retry_at <= ? AND retry_count < max_retries", domain.MintRequestStatusFailed, now).
		Order("next_retry_at ASC").
		Limit(limit).
		Find(&requests).Error
	return requests, err
}

func (r Repository) MarkCompleted(ctx context.Context, id uuid.UUID, txHash string, blockNumber int64) error {
	return r.db.WithContext(ctx).Model(&domain.MintRequest{}).Where("id = ?", id).Updates(map[string]any{
		"status":           domain.MintRequestStatusCompleted,
		"transaction_hash": txHash,
		"block_number":     blockNumber,
		"completed_at":     new(time.Now().UTC()),
	}).Error
}

func (r Repository) MarkFailed(ctx context.Context, id uuid.UUID, errorCode, errorMessage string, nextRetryAt *time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.MintRequest{}).Where("id = ?", id).Updates(map[string]any{
		"status":        domain.MintRequestStatusFailed,
		"error_code":    errorCode,
		"error_message": errorMessage,
		"next_retry_at": nextRetryAt,
		"retry_count":   gorm.Expr("retry_count + 1"),
	}).Error
}
