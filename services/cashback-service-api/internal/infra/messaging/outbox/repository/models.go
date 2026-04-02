package repository

import (
	"time"

	outboxdomain "github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox/domain"
)

type outboxModel struct {
	ID            int64               `gorm:"primaryKey;autoIncrement"`
	EventType     string              `gorm:"not null"`
	AggregateType string              `gorm:"not null"`
	AggregateID   int64               `gorm:"not null"`
	Payload       []byte              `gorm:"not null"`
	Status        outboxdomain.Status `gorm:"not null;default:'pending'"`
	RetryCount    int                 `gorm:"default:0"`
	MaxRetries    int                 `gorm:"default:5"`
	ErrorMessage  string
	CreatedAt     time.Time
	PublishedAt   *time.Time
	UpdatedAt     time.Time
}

func (outboxModel) TableName() string {
	return "cashback.outbox_events"
}

func (m outboxModel) toDomain() outboxdomain.Event {
	return outboxdomain.Event{
		ID:            m.ID,
		EventType:     m.EventType,
		AggregateType: m.AggregateType,
		AggregateID:   m.AggregateID,
		Payload:       m.Payload,
		Status:        m.Status,
		RetryCount:    m.RetryCount,
		MaxRetries:    m.MaxRetries,
		ErrorMessage:  m.ErrorMessage,
	}
}
