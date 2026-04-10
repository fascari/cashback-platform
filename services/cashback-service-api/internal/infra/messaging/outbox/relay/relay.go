package relay

import (
	"context"
	"log"
	"time"

	outboxdomain "github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox/domain"
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/messaging/outbox/repository"
	"github.com/cashback-platform/services/cashback-service-api/internal/infra/nats"
	"go.uber.org/fx"
)

type Relay struct {
	repo     repository.Repository
	nats     *nats.NATSClient
	interval time.Duration
	done     chan struct{}
}

func New(repo repository.Repository, natsClient *nats.NATSClient, interval time.Duration) Relay {
	return Relay{
		repo:     repo,
		nats:     natsClient,
		interval: interval,
		done:     make(chan struct{}),
	}
}

func (r Relay) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-r.done:
			return
		case <-ticker.C:
			r.poll(ctx)
		}
	}
}

func (r Relay) Stop() {
	close(r.done)
}

func (r Relay) poll(ctx context.Context) {
	events, err := r.repo.Pending(ctx, 100)
	if err != nil {
		log.Printf("outbox relay: error fetching pending events: %v", err)
		return
	}

	for _, event := range events {
		r.dispatch(ctx, event)
	}
}

func (r Relay) dispatch(ctx context.Context, event outboxdomain.Event) {
	if err := r.nats.Publish(event.EventType, event.Payload); err != nil {
		r.handleError(ctx, event, err)
		return
	}

	if err := r.repo.MarkAsPublished(ctx, event.ID); err != nil {
		log.Printf("outbox relay: error marking event %d as published: %v", event.ID, err)
	}
}

func (r Relay) handleError(ctx context.Context, event outboxdomain.Event, publishErr error) {
	log.Printf("outbox relay: error dispatching event %d: %v", event.ID, publishErr)

	if err := r.repo.IncrementRetry(ctx, event.ID); err != nil {
		log.Printf("outbox relay: error incrementing retry for event %d: %v", event.ID, err)
	}

	if event.RetryCount >= event.MaxRetries-1 {
		if err := r.repo.MarkAsFailed(ctx, event.ID, publishErr.Error()); err != nil {
			log.Printf("outbox relay: error marking event %d as failed: %v", event.ID, err)
		}
	}
}

func Start(lc fx.Lifecycle, r Relay) {
	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go r.Run(ctx)
			log.Println("outbox relay: started")
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			r.Stop()
			log.Println("outbox relay: stopped")
			return nil
		},
	})
}
