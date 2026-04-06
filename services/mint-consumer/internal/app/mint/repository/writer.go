package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/domain"
	"gorm.io/gorm"
)

func (r Repository) CreateMintRequest(ctx context.Context, req domain.MintRequest) (domain.MintRequest, error) {
	m := new(mintRequestFromDomain(req))
	if err := r.conn(ctx).Create(m).Error; err != nil {
		return domain.MintRequest{}, err
	}
	return m.toDomain(), nil
}

func (r Repository) MarkCompleted(ctx context.Context, id int64, txHash string, blockNumber int64) error {
	return r.conn(ctx).Model(&mintRequestModel{}).Where("id = ?", id).Updates(map[string]any{
		"status":           domain.MintRequestStatusCompleted,
		"transaction_hash": txHash,
		"block_number":     blockNumber,
		"completed_at":     new(time.Now().UTC()),
	}).Error
}

func (r Repository) MarkFailed(ctx context.Context, id int64, errorCode, errorMessage string, nextRetryAt *time.Time) error {
	return r.conn(ctx).Model(&mintRequestModel{}).Where("id = ?", id).Updates(map[string]any{
		"status":        domain.MintRequestStatusFailed,
		"error_code":    errorCode,
		"error_message": errorMessage,
		"next_retry_at": nextRetryAt,
		"retry_count":   gorm.Expr("retry_count + 1"),
	}).Error
}

func (r Repository) CreateProcessedEvent(ctx context.Context, eventID uuid.UUID, eventType string) error {
	m := new(processedEventModel{EventID: eventID, EventType: eventType})
	return r.conn(ctx).Create(m).Error
}
