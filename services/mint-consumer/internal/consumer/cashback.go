package consumer

import (
	"context"
	"errors"
	"time"

	"github.com/cashback-platform/kit/logger"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/processcashbackapproved"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/retryfailedmints"
	"github.com/cashback-platform/services/mint-consumer/internal/infra/nats"
	natsgo "github.com/nats-io/nats.go"
	"go.uber.org/fx"
)

// CashbackConsumer subscribes to cashback.approved events and drives mint operations.
type CashbackConsumer struct {
	processUseCase processcashbackapproved.UseCase
	retryUseCase   retryfailedmints.UseCase
	natsClient     *nats.NATSClient
	done           chan struct{}
	sub            *natsgo.Subscription
}

// NewCashback creates a CashbackConsumer wired with its use cases and NATS client.
func NewCashback(
	processUseCase processcashbackapproved.UseCase,
	retryUseCase retryfailedmints.UseCase,
	natsClient *nats.NATSClient,
) *CashbackConsumer {
	return &CashbackConsumer{
		processUseCase: processUseCase,
		retryUseCase:   retryUseCase,
		natsClient:     natsClient,
		done:           make(chan struct{}),
	}
}

// Start subscribes to cashback.approved and begins processing and retry loops.
func (c *CashbackConsumer) Start(ctx context.Context) error {
	js := c.natsClient.JetStream()

	consumerConfig := &natsgo.ConsumerConfig{
		Durable:       "mint-consumer",
		FilterSubject: "cashback.approved",
		DeliverPolicy: natsgo.DeliverAllPolicy,
		AckPolicy:     natsgo.AckExplicitPolicy,
		MaxDeliver:    5,
		AckWait:       30 * time.Second,
	}

	_, err := js.AddConsumer("CASHBACK_EVENTS", consumerConfig)
	if err != nil && !errors.Is(err, natsgo.ErrConsumerNameAlreadyInUse) {
		logger.Warn("failed to create NATS consumer", "error", err)
	}

	sub, err := js.PullSubscribe("cashback.approved", "mint-consumer")
	if err != nil {
		return err
	}
	c.sub = sub

	logger.Info("cashback consumer started")

	go c.processMessages(ctx)
	go c.retryLoop(ctx)

	return nil
}

func (c *CashbackConsumer) processMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		default:
			msgs, err := c.sub.Fetch(10, natsgo.MaxWait(time.Second))
			if err != nil {
				if !errors.Is(err, natsgo.ErrTimeout) {
					logger.Error("error fetching messages", "error", err)
				}
				continue
			}
			for _, msg := range msgs {
				c.handleMessage(ctx, msg)
			}
		}
	}
}

func (c *CashbackConsumer) handleMessage(ctx context.Context, msg *natsgo.Msg) {
	if err := c.processUseCase.Execute(ctx, msg.Data); err != nil {
		logger.Error("error processing cashback approved message", "error", err)
		if nakErr := msg.Nak(); nakErr != nil {
			logger.Error("error NAKing message", "error", nakErr)
		}
		return
	}
	if ackErr := msg.Ack(); ackErr != nil {
		logger.Error("error ACKing message", "error", ackErr)
	}
}

func (c *CashbackConsumer) retryLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		case <-ticker.C:
			if err := c.retryUseCase.Execute(ctx); err != nil {
				logger.Error("error retrying failed mints", "error", err)
			}
		}
	}
}

// Stop closes the done channel and unsubscribes from NATS.
func (c *CashbackConsumer) Stop() {
	close(c.done)
	if c.sub != nil {
		if err := c.sub.Unsubscribe(); err != nil {
			logger.Error("error unsubscribing from NATS", "error", err)
		}
	}
}

// StartConsumer registers Start/Stop hooks with the fx lifecycle.
func StartConsumer(lc fx.Lifecycle, c *CashbackConsumer) {
	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if err := c.Start(ctx); err != nil {
				return err
			}
			logger.Info("cashback consumer started via fx lifecycle")
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			c.Stop()
			logger.Info("cashback consumer stopped")
			return nil
		},
	})
}
