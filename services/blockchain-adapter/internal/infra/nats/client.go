package nats

import (
	"fmt"

	natsgo "github.com/nats-io/nats.go"

	"github.com/cashback-platform/services/blockchain-adapter/internal/config"
)

// NATSClient wraps a NATS connection with JetStream publishing support.
// Streams must exist before this client is created — run cmd/nats-setup during infra setup.
type NATSClient struct {
	conn *natsgo.Conn
	js   natsgo.JetStreamContext
}

func NewNATSClient(cfg *config.Config) (*NATSClient, error) {
	conn, err := natsgo.Connect(cfg.NATS.URL)
	if err != nil {
		return nil, fmt.Errorf("connect to NATS: %w", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("create JetStream context: %w", err)
	}

	return new(NATSClient{conn: conn, js: js}), nil
}

// Publish sends a message to the given JetStream subject.
func (c *NATSClient) Publish(subject string, data []byte) error {
	_, err := c.js.Publish(subject, data)
	return err
}

// Close drains and closes the underlying NATS connection.
func (c *NATSClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
