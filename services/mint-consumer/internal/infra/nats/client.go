package nats

import (
	"fmt"
	"log"

	"github.com/cashback-platform/services/mint-consumer/internal/config"
	"github.com/nats-io/nats.go"
)

type NATSClient struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func NewNATSClient(cfg *config.Config) (*NATSClient, error) {
	conn, err := nats.Connect(cfg.NATS.URL)
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

func (c *NATSClient) JetStream() nats.JetStreamContext {
	return c.js
}

func (c *NATSClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
