package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Repository struct {
		db *gorm.DB
	}

	OutboxEvent struct {
		ID         uuid.UUID
		EventType  string
		Payload    []byte
		RetryCount int
		MaxRetries int
		Published  bool
		Failed     bool
		Error      string
	}
)

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, eventType string, payload []byte) error {
	event := outboxModel{
		ID:         uuid.New(),
		EventType:  eventType,
		Payload:    payload,
		RetryCount: 0,
		MaxRetries: 3,
		Published:  false,
		Failed:     false,
	}
	return r.db.WithContext(ctx).Create(&event).Error
}

func (r *Repository) Pending(ctx context.Context, limit int) ([]OutboxEvent, error) {
	var models []outboxModel
	if err := r.db.WithContext(ctx).
		Where("published = ? AND failed = ?", false, false).
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	events := make([]OutboxEvent, len(models))
	for i, m := range models {
		d := toDomain(&m)
		events[i] = OutboxEvent{
			ID:         d.ID,
			EventType:  d.EventType,
			Payload:    d.Payload,
			RetryCount: d.RetryCount,
			MaxRetries: d.MaxRetries,
			Published:  d.Published,
			Failed:     d.Failed,
			Error:      d.Error,
		}
	}
	return events, nil
}

func (r *Repository) IncrementRetry(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&outboxModel{}).
		Where("id = ?", id).
		Update("retry_count", gorm.Expr("retry_count + ?", 1)).Error
}

func (r *Repository) MarkAsPublished(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&outboxModel{}).
		Where("id = ?", id).
		Update("published", true).Error
}

func (r *Repository) MarkAsFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	return r.db.WithContext(ctx).
		Model(&outboxModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"failed": true,
			"error":  errMsg,
		}).Error
}
