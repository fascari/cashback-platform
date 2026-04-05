package nats

import (
	"fmt"

	kitnats "github.com/cashback-platform/kit/nats"
	"github.com/cashback-platform/services/mint-consumer/internal/config"
	natsgo "github.com/nats-io/nats.go"
)

// NATSClient wraps kit/nats.Client for mint-consumer.
type NATSClient struct {
	*kitnats.Client
}

func NewClient(cfg config.NATS) (*NATSClient, error) {
	c, err := kitnats.New(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("nats client: %w", err)
	}
	return &NATSClient{Client: c}, nil
}

// JetStream exposes the underlying JetStreamContext for consumer setup.
func (c *NATSClient) JetStream() natsgo.JetStreamContext {
	return c.Client.JetStream()
}

