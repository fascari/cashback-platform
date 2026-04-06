package repository

import (
	"context"
	"time"

	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/domain"
	"github.com/google/uuid"
)

func (r Repository) FindFailedRetryable(ctx context.Context, limit int) ([]domain.MintRequest, error) {
	var models []mintRequestModel
	now := time.Now().UTC()
	err := r.conn(ctx).
		Where("status = ? AND next_retry_at <= ? AND retry_count < max_retries",
			domain.MintRequestStatusFailed, now).
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

func (r Repository) ExistsProcessedEvent(ctx context.Context, eventID uuid.UUID) (bool, error) {
	count := new(int64)
	err := r.conn(ctx).Model(&processedEventModel{}).Where("event_id = ?", eventID).Count(count).Error
	return *count > 0, err
}
