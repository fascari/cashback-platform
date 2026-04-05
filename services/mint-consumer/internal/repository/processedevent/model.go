package processedevent

import (
	"time"

	"github.com/google/uuid"

	"github.com/cashback-platform/services/mint-consumer/internal/domain"
)

type processedEventModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	EventID   uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	EventType string    `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (processedEventModel) TableName() string {
	return "mint.processed_events"
}

func (m processedEventModel) toDomain() domain.ProcessedEvent {
	return domain.ProcessedEvent{
		ID:        m.ID,
		EventID:   m.EventID,
		EventType: m.EventType,
		CreatedAt: m.CreatedAt,
	}
}

func fromDomain(e domain.ProcessedEvent) processedEventModel {
	return processedEventModel{
		ID:        e.ID,
		EventID:   e.EventID,
		EventType: e.EventType,
	}
}
