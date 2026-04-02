package repository

type outboxEvent struct {
	ID            int64
	EventType     string
	AggregateType string
	AggregateID   int64
	Payload       []byte
	Status        string
	RetryCount    int
	MaxRetries    int
	ErrorMessage  string
}

func toDomain(m *outboxModel) *outboxEvent {
	return &outboxEvent{
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
