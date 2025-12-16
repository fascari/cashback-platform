package repository

import (
	"github.com/google/uuid"
)

type outboxEvent struct {
	ID         uuid.UUID
	EventType  string
	Payload    []byte
	RetryCount int
	MaxRetries int
	Published  bool
	Failed     bool
	Error      string
}

func toDomain(m *outboxModel) *outboxEvent {
	return &outboxEvent{
		ID:         m.ID,
		EventType:  m.EventType,
		Payload:    m.Payload,
		RetryCount: m.RetryCount,
		MaxRetries: m.MaxRetries,
		Published:  m.Published,
		Failed:     m.Failed,
		Error:      m.Error,
	}
}
