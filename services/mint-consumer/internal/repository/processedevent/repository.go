package processedevent

import (
	"context"

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

func (r Repository) Create(ctx context.Context, event domain.ProcessedEvent) error {
	m := new(fromDomain(event))
	db := r.db.WithContext(ctx)
	if tx := gormtx.ExtractTx(ctx); tx != nil {
		db = tx.WithContext(ctx)
	}
	return db.Create(m).Error
}

func (r Repository) Exists(ctx context.Context, eventID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&processedEventModel{}).Where("event_id = ?", eventID).Count(&count).Error
	return count > 0, err
}

func (r Repository) FindByEventID(ctx context.Context, eventID uuid.UUID) (domain.ProcessedEvent, error) {
	m := new(processedEventModel)
	if err := r.db.WithContext(ctx).Where("event_id = ?", eventID).First(m).Error; err != nil {
		return domain.ProcessedEvent{}, err
	}
	return m.toDomain(), nil
}
