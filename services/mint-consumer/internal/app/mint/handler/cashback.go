package handler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cashback-platform/kit/events"
	"github.com/cashback-platform/kit/logger"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/handler/cashbackapproved"
	"github.com/cashback-platform/services/mint-consumer/internal/app/mint/usecase/retrymints"
	"github.com/cashback-platform/services/mint-consumer/internal/infra/nats"
	natsgo "github.com/nats-io/nats.go"
	"go.uber.org/fx"
)

type CashbackConsumer struct {
	approvedHandler cashbackapproved.Handler
	retryUseCase    retrymints.UseCase
	natsClient      *nats.NATSClient
	sub             *natsgo.Subscription
	wg              sync.WaitGroup
}

func NewCashback(
	approvedHandler cashbackapproved.Handler,
	retryUseCase retrymints.UseCase,
	natsClient *nats.NATSClient,
) *CashbackConsumer {
	return &CashbackConsumer{
		approvedHandler: approvedHandler,
		retryUseCase:    retryUseCase,
		natsClient:      natsClient,
	}
}

func (c *CashbackConsumer) start(ctx context.Context) error {
	js := c.natsClient.JetStream()

	if _, err := js.AddConsumer(events.StreamCashbackEvents, &natsgo.ConsumerConfig{
		Durable:       "mint-consumer",
		FilterSubject: events.CashbackApproved,
		DeliverPolicy: natsgo.DeliverAllPolicy,
		AckPolicy:     natsgo.AckExplicitPolicy,
		MaxDeliver:    5,
		AckWait:       30 * time.Second,
	}); err != nil && !errors.Is(err, natsgo.ErrConsumerNameAlreadyInUse) {
		return fmt.Errorf("add NATS consumer: %w", err)
	}

	var err error
	c.sub, err = js.PullSubscribe(events.CashbackApproved, "mint-consumer")
	if err != nil {
		return err
	}

	c.safeGo(func() { c.processMessages(ctx) })
	c.safeGo(func() { c.retryLoop(ctx) })

	logger.Info("cashback consumer started")
	return nil
}

// safeGo runs fn in a goroutine tracked by wg, recovering from panics to prevent silent goroutine death.
func (c *CashbackConsumer) safeGo(fn func()) {
	c.wg.Go(func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic in consumer goroutine", "recover", r)
			}
		}()
		fn()
	})
}

func (c *CashbackConsumer) stop() {
	c.wg.Wait()
	if c.sub != nil {
		if err := c.sub.Unsubscribe(); err != nil {
			logger.Error("error unsubscribing from NATS", "error", err)
		}
	}
}

func (c *CashbackConsumer) processMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msgs, err := c.sub.Fetch(10, natsgo.MaxWait(time.Second))
			if err != nil {
				if !errors.Is(err, natsgo.ErrTimeout) {
					logger.Error("error fetching messages", "error", err)
				}
				continue
			}
			c.dispatch(ctx, msgs)
		}
	}
}

func (c *CashbackConsumer) dispatch(ctx context.Context, msgs []*natsgo.Msg) {
	for _, msg := range msgs {
		if err := c.approvedHandler.Handle(ctx, msg); err != nil {
			logger.Error("error handling cashback.approved message", "error", err)
			nak(msg)
			continue
		}
		if err := msg.Ack(); err != nil {
			logger.Error("error ACKing message", "error", err)
		}
	}
}

func (c *CashbackConsumer) retryLoop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.retryUseCase.Execute(ctx); err != nil {
				logger.Error("error retrying failed mints", "error", err)
			}
		}
	}
}

func nak(msg *natsgo.Msg) {
	if err := msg.Nak(); err != nil {
		logger.Error("error NAKing message", "error", err)
	}
}

func StartConsumer(lc fx.Lifecycle, c *CashbackConsumer) {
	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if err := c.start(ctx); err != nil {
				cancel()
				return err
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			c.stop()
			logger.Info("cashback consumer stopped")
			return nil
		},
	})
}
