package consumer

import (
	"context"
	"log"
	"time"

	"github.com/cashback-platform/services/mint-consumer/internal/infra/nats"
	"github.com/cashback-platform/services/mint-consumer/internal/usecase"
	natsgo "github.com/nats-io/nats.go"
	"go.uber.org/fx"
)

type CashbackConsumer struct {
	mintUsecase *usecase.MintUsecase
	natsClient  *nats.NATSClient
	done        chan struct{}
	sub         *natsgo.Subscription
}

func NewCashbackConsumer(mintUsecase *usecase.MintUsecase, natsClient *nats.NATSClient) *CashbackConsumer {
	return &CashbackConsumer{
		mintUsecase: mintUsecase,
		natsClient:  natsClient,
		done:        make(chan struct{}),
	}
}

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
	if err != nil && err != natsgo.ErrConsumerNameAlreadyInUse {
		log.Printf("Warning: Failed to create consumer: %v", err)
	}

	sub, err := js.PullSubscribe("cashback.approved", "mint-consumer")
	if err != nil {
		return err
	}
	c.sub = sub

	log.Println("Cashback consumer started, listening for cashback.approved events")

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
				if err != natsgo.ErrTimeout {
					log.Printf("Error fetching messages: %v", err)
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
	log.Printf("Processing message: %s", string(msg.Data))

	if err := c.mintUsecase.ProcessCashbackApproved(ctx, msg.Data); err != nil {
		log.Printf("Error processing message: %v", err)
		if err := msg.Nak(); err != nil {
			log.Printf("Error NAKing message: %v", err)
		}
	} else {
		if err := msg.Ack(); err != nil {
			log.Printf("Error ACKing message: %v", err)
		}
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
			if err := c.mintUsecase.RetryFailedMints(ctx); err != nil {
				log.Printf("Error retrying failed mints: %v", err)
			}
		}
	}
}

func (c *CashbackConsumer) Stop() {
	close(c.done)
	if c.sub != nil {
		if err := c.sub.Unsubscribe(); err != nil {
			log.Printf("Error unsubscribing: %v", err)
		}
	}
}

func StartConsumer(lc fx.Lifecycle, consumer *CashbackConsumer) {
	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if err := consumer.Start(ctx); err != nil {
				return err
			}
			log.Println("Cashback consumer started")
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			consumer.Stop()
			log.Println("Cashback consumer stopped")
			return nil
		},
	})
}
