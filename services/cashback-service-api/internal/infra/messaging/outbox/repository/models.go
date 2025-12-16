package repository

import (
	"time"

	"github.com/google/uuid"
)

type outboxModel struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key"`
	EventType  string    `gorm:"not null"`
	Payload    []byte    `gorm:"not null"`
	RetryCount int       `gorm:"default:0"`
	MaxRetries int       `gorm:"default:3"`
	Published  bool      `gorm:"default:false"`
	Failed     bool      `gorm:"default:false"`
	Error      string
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (outboxModel) TableName() string {
	return "outbox_events"
}
