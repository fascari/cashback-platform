package outbox

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox/repository"
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/nats"
	"go.uber.org/fx"
)

// OutboxPublisher implements the EventPublisher interface using the Outbox Pattern
// It persists events before publishing to ensure reliable delivery
type OutboxPublisher struct {
	outboxRepo *repository.Repository
	natsClient *nats.NATSClient
	done       chan struct{}
}

func NewOutboxPublisher(outboxRepo *repository.Repository, natsClient *nats.NATSClient) *OutboxPublisher {
	return &OutboxPublisher{
		outboxRepo: outboxRepo,
		natsClient: natsClient,
		done:       make(chan struct{}),
	}
}

// Publish adds an event to the outbox for async publishing
// Implements messaging.EventPublisher interface
func (p *OutboxPublisher) Publish(ctx context.Context, eventType string, payload any) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.outboxRepo.Create(ctx, eventType, payloadBytes)
}

func (p *OutboxPublisher) Start(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.done:
			return
		case <-ticker.C:
			p.processEvents(ctx)
		}
	}
}

func (p *OutboxPublisher) Stop() {
	close(p.done)
}

func (p *OutboxPublisher) processEvents(ctx context.Context) {
	events, err := p.outboxRepo.Pending(ctx, 100)
	if err != nil {
		log.Printf("Error fetching pending events: %v", err)
		return
	}

	for _, event := range events {
		p.publishEvent(ctx, event)
	}
}

func (p *OutboxPublisher) publishEvent(ctx context.Context, event repository.OutboxEvent) {
	subject := event.EventType

	if err := p.natsClient.Publish(subject, event.Payload); err != nil {
		p.handlePublishError(ctx, event, err)
		return
	}

	if err := p.outboxRepo.MarkAsPublished(ctx, event.ID); err != nil {
		log.Printf("Error marking event %s as published: %v", event.ID, err)
	}
}

func (p *OutboxPublisher) handlePublishError(ctx context.Context, event repository.OutboxEvent, publishErr error) {
	log.Printf("Error publishing event %s: %v", event.ID, publishErr)

	if err := p.outboxRepo.IncrementRetry(ctx, event.ID); err != nil {
		log.Printf("Error incrementing retry for event %s: %v", event.ID, err)
	}

	if event.RetryCount >= event.MaxRetries-1 {
		if err := p.outboxRepo.MarkAsFailed(ctx, event.ID, publishErr.Error()); err != nil {
			log.Printf("Error marking event %s as failed: %v", event.ID, err)
		}
	}
}

func StartOutboxPublisher(lc fx.Lifecycle, publisher *OutboxPublisher) {
	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go publisher.Start(ctx)
			log.Println("Outbox publisher started")
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			publisher.Stop()
			log.Println("Outbox publisher stopped")
			return nil
		},
	})
}
