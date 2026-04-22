package handler

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/cashback-platform/kit/events"
	"github.com/cashback-platform/kit/logger"
	"github.com/cashback-platform/services/cashback-service-api/internal/app/deposit/handler/depositdetected"
	infranats "github.com/cashback-platform/services/cashback-service-api/internal/infra/nats"
	natsgo "github.com/nats-io/nats.go"
	"go.uber.org/fx"
)

const (
	consumerName = "cashback-service-api-deposits"
	fetchBatch   = 10
	fetchTimeout = time.Second
)

type DepositConsumer struct {
	js      natsgo.JetStreamContext
	handler depositdetected.Handler
	sub     *natsgo.Subscription
	wg      sync.WaitGroup
}

func NewDepositConsumer(client *infranats.NATSClient, h depositdetected.Handler) *DepositConsumer {
	return &DepositConsumer{
		js:      client.JetStream(),
		handler: h,
	}
}

func (c *DepositConsumer) start(ctx context.Context) error {
	var err error
	c.sub, err = c.js.PullSubscribe(events.DepositDetected, consumerName,
		natsgo.BindStream(events.StreamDepositEvents),
	)
	if err != nil {
		return err
	}
	c.safeGo(func() { c.processMessages(ctx) })
	logger.Info("deposit consumer started")
	return nil
}

func (c *DepositConsumer) safeGo(fn func()) {
	c.wg.Go(func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic in deposit consumer goroutine", "recover", r)
			}
		}()
		fn()
	})
}

func (c *DepositConsumer) stop() {
	c.wg.Wait()
	if c.sub != nil {
		if err := c.sub.Unsubscribe(); err != nil {
			logger.Error("error unsubscribing deposit consumer", "error", err)
		}
	}
}

func (c *DepositConsumer) processMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msgs, err := c.sub.Fetch(fetchBatch, natsgo.MaxWait(fetchTimeout))
		if err != nil {
			if !errors.Is(err, natsgo.ErrTimeout) {
				logger.Error("error fetching deposit messages", "error", err)
			}
			continue
		}

		for _, msg := range msgs {
			if handleErr := c.handler.Handle(ctx, msg.Data); handleErr != nil {
				logger.Error("error handling deposit.detected message", "error", handleErr)
				if depositdetected.IsPermanent(handleErr) {
					if termErr := msg.Term(); termErr != nil {
						logger.Error("error terming deposit message", "error", termErr)
					}
				} else {
					if nakErr := msg.Nak(); nakErr != nil {
						logger.Error("error naking deposit message", "error", nakErr)
					}
				}
				continue
			}
			if ackErr := msg.Ack(); ackErr != nil {
				logger.Error("error acking deposit message", "error", ackErr)
			}
		}
	}
}

func StartConsumer(lc fx.Lifecycle, consumer *DepositConsumer) {
	ctx, cancel := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if err := consumer.start(ctx); err != nil {
				cancel()
				return err
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			consumer.stop()
			logger.Info("deposit consumer stopped")
			return nil
		},
	})
}
