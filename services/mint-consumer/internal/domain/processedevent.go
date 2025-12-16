package domain

import (
	"time"

	"github.com/google/uuid"
)

// ProcessedEvent tracks events that have been processed for idempotency
type ProcessedEvent struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EventID     uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	EventType   string    `gorm:"type:varchar(100);not null"`
	ProcessedAt time.Time `gorm:"autoCreateTime"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

// TableName specifies the table name for GORM
func (ProcessedEvent) TableName() string {
	return "processed_events"
}
