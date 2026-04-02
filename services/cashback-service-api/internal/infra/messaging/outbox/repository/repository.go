package repository

import (
	"context"
	"time"

	outboxdomain "github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox/domain"
	"gorm.io/gorm"
)

type Repository struct {
	db         *gorm.DB
	maxRetries int
}

func New(db *gorm.DB, maxRetries int) Repository {
	return Repository{db: db, maxRetries: maxRetries}
}

func (r Repository) Create(ctx context.Context, eventType, aggregateType string, aggregateID int64, payload []byte) error {
	event := outboxModel{
		EventType:     eventType,
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		Payload:       payload,
		Status:        outboxdomain.StatusPending,
		MaxRetries:    r.maxRetries,
	}
	return r.db.WithContext(ctx).Create(&event).Error
}

func (r Repository) Pending(ctx context.Context, limit int) ([]outboxdomain.Event, error) {
	var models []outboxModel
	if err := r.db.WithContext(ctx).
		Where("status = ?", outboxdomain.StatusPending).
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	events := make([]outboxdomain.Event, len(models))
	for i, m := range models {
		events[i] = m.toDomain()
	}
	return events, nil
}

func (r Repository) IncrementRetry(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&outboxModel{}).
		Where("id = ?", id).
		Update("retry_count", gorm.Expr("retry_count + ?", 1)).Error
}

func (r Repository) MarkAsPublished(ctx context.Context, id int64) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).
		Model(&outboxModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       string(outboxdomain.StatusPublished),
			"published_at": now,
		}).Error
}

func (r Repository) MarkAsFailed(ctx context.Context, id int64, errMsg string) error {
	return r.db.WithContext(ctx).
		Model(&outboxModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":        string(outboxdomain.StatusFailed),
			"error_message": errMsg,
		}).Error
}
