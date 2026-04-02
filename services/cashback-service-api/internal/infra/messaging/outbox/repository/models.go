package repository

import "time"

type outboxModel struct {
	ID            int64      `gorm:"primaryKey;autoIncrement"`
	EventType     string     `gorm:"not null"`
	AggregateType string     `gorm:"not null"`
	AggregateID   int64      `gorm:"not null"`
	Payload       []byte     `gorm:"not null"`
	Status        string     `gorm:"not null;default:'pending'"`
	RetryCount    int        `gorm:"default:0"`
	MaxRetries    int        `gorm:"default:5"`
	ErrorMessage  string
	CreatedAt     time.Time
	PublishedAt   *time.Time
	UpdatedAt     time.Time
}

func (outboxModel) TableName() string { return "outbox_events" }
