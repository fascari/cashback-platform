package domain

type Status string

const (
	StatusPending   Status = "pending"
	StatusPublished Status = "published"
	StatusFailed    Status = "failed"
)

type Event struct {
	ID            int64
	EventType     string
	AggregateType string
	AggregateID   int64
	Payload       []byte
	Status        Status
	RetryCount    int
	MaxRetries    int
	ErrorMessage  string
}
