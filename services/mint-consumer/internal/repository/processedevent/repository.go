package processedevent

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/cashback-platform/services/mint-consumer/internal/domain"
)

// Repository provides data access for processed events.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new processed event repository.
func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

// Create records an event as processed.
func (r Repository) Create(ctx context.Context, event domain.ProcessedEvent) error {
	m := fromDomain(event)
	return r.db.WithContext(ctx).Create(&m).Error
}

// Exists reports whether an event with the given ID has already been processed.
func (r Repository) Exists(ctx context.Context, eventID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&processedEventModel{}).Where("event_id = ?", eventID).Count(&count).Error
	return count > 0, err
}

// FindByEventID retrieves a processed event record by its event UUID.
func (r Repository) FindByEventID(ctx context.Context, eventID uuid.UUID) (domain.ProcessedEvent, error) {
	m := new(processedEventModel)
	if err := r.db.WithContext(ctx).Where("event_id = ?", eventID).First(m).Error; err != nil {
		return domain.ProcessedEvent{}, err
	}
	return m.toDomain(), nil
}
