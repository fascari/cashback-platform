package repository

import (
	"context"
	"time"

	"github.com/cashback-platform/services/mint-consumer/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	ProcessedEventRepository interface {
		Create(ctx context.Context, event *domain.ProcessedEvent) error
		Exists(ctx context.Context, eventID uuid.UUID) (bool, error)
		GetByEventID(ctx context.Context, eventID uuid.UUID) (*domain.ProcessedEvent, error)
	}

	processedEventRepository struct {
		db *gorm.DB
	}
)

func NewProcessedEventRepository(db *gorm.DB) ProcessedEventRepository {
	return &processedEventRepository{db: db}
}

func (r *processedEventRepository) Create(ctx context.Context, event *domain.ProcessedEvent) error {
	event.ProcessedAt = time.Now().UTC()
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *processedEventRepository) Exists(ctx context.Context, eventID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.ProcessedEvent{}).Where("event_id = ?", eventID).Count(&count).Error
	return count > 0, err
}

func (r *processedEventRepository) GetByEventID(ctx context.Context, eventID uuid.UUID) (*domain.ProcessedEvent, error) {
	var event domain.ProcessedEvent
	if err := r.db.WithContext(ctx).Where("event_id = ?", eventID).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}
