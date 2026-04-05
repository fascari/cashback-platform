package processedevent

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

func (r Repository) Create(ctx context.Context, event *domain.ProcessedEvent) error {
	event.ProcessedAt = time.Now().UTC()
	return r.db.WithContext(ctx).Create(event).Error
}

func (r Repository) Exists(ctx context.Context, eventID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.ProcessedEvent{}).Where("event_id = ?", eventID).Count(&count).Error
	return count > 0, err
}

func (r Repository) FindByEventID(ctx context.Context, eventID uuid.UUID) (*domain.ProcessedEvent, error) {
	event := new(domain.ProcessedEvent)
	if err := r.db.WithContext(ctx).Where("event_id = ?", eventID).First(event).Error; err != nil {
		return nil, err
	}
	return event, nil
}
