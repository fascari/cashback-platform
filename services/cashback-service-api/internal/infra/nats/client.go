package nats

import (
	"fmt"
	"log"

	"github.com/cashback-platform/services/cashback-service-api/internal/config"
	natsgo "github.com/nats-io/nats.go"
)

// NATSClient wraps a NATS connection with JetStream publishing support.
// Streams must exist before this client is created — run cmd/nats-setup during infra setup.
type NATSClient struct {
	conn *natsgo.Conn
	js   natsgo.JetStreamContext
}

func NewNATSClient(cfg config.NATS) (*NATSClient, error) {
	conn, err := natsgo.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	log.Println("NATS connected successfully")
	return &NATSClient{
		conn: conn,
		js:   js,
	}, nil
}

func (c *NATSClient) Publish(subject string, data []byte) error {
	_, err := c.js.Publish(subject, data)
	return err
}

func (c *NATSClient) JetStream() natsgo.JetStreamContext {
	return c.js
}

func (c *NATSClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
