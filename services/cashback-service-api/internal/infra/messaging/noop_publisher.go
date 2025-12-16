package messaging

import "context"

// NoopPublisher is an EventPublisher that performs no operations.
// It is useful for testing or when event publishing needs to be disabled.
type NoopPublisher struct{}

func NewNoopPublisher() *NoopPublisher {
	return &NoopPublisher{}
}

func (*NoopPublisher) Publish(_ context.Context, _ string, _ any) error {
	return nil
}
