package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProcessedEvent struct {
	ID        int64
	EventID   uuid.UUID
	EventType string
	CreatedAt time.Time
}
