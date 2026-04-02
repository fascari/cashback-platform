package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type (
	Repository struct {
		db *gorm.DB
	}

	OutboxEvent struct {
		ID            int64
		EventType     string
		AggregateType string
		AggregateID   int64
		Payload       []byte
		Status        string
		RetryCount    int
		MaxRetries    int
		ErrorMessage  string
	}
)

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, eventType, aggregateType string, aggregateID int64, payload []byte) error {
	event := outboxModel{
		EventType:     eventType,
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		Payload:       payload,
		Status:        "pending",
		RetryCount:    0,
		MaxRetries:    5,
	}
	return r.db.WithContext(ctx).Create(&event).Error
}

func (r *Repository) Pending(ctx context.Context, limit int) ([]OutboxEvent, error) {
	var models []outboxModel
	if err := r.db.WithContext(ctx).
		Where("status = ?", "pending").
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	events := make([]OutboxEvent, len(models))
	for i, m := range models {
		d := toDomain(&m)
		events[i] = OutboxEvent{
			ID:            d.ID,
			EventType:     d.EventType,
			AggregateType: d.AggregateType,
			AggregateID:   d.AggregateID,
			Payload:       d.Payload,
			Status:        d.Status,
			RetryCount:    d.RetryCount,
			MaxRetries:    d.MaxRetries,
			ErrorMessage:  d.ErrorMessage,
		}
	}
	return events, nil
}

func (r *Repository) IncrementRetry(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&outboxModel{}).
		Where("id = ?", id).
		Update("retry_count", gorm.Expr("retry_count + ?", 1)).Error
}

func (r *Repository) MarkAsPublished(ctx context.Context, id int64) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).
		Model(&outboxModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       "published",
			"published_at": now,
		}).Error
}

func (r *Repository) MarkAsFailed(ctx context.Context, id int64, errMsg string) error {
	return r.db.WithContext(ctx).
		Model(&outboxModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":        "failed",
			"error_message": errMsg,
		}).Error
}
