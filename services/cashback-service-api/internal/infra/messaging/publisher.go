// Package messaging provides abstractions for asynchronous event publishing.
// It defines interfaces that can be implemented by different messaging strategies
// such as Outbox Pattern, direct NATS, Kafka, etc.
package messaging

import "context"

// EventPublisher publishes domain events asynchronously.
// Implementations must handle serialization, delivery, and error handling.
type EventPublisher interface {
	// Publish publishes an event of the given type with the specified payload.
	// The actual publishing mechanism is implementation-specific.
	// Returns an error if the event cannot be published or persisted.
	Publish(ctx context.Context, eventType string, payload any) error
}
