package nats

import (
	"fmt"

	natsgo "github.com/nats-io/nats.go"
)

// Client wraps a NATS JetStream connection.
type Client struct {
	conn *natsgo.Conn
	js   natsgo.JetStreamContext
}

// New connects to NATS at url and returns a JetStream-enabled client.
func New(url string) (*Client, error) {
	conn, err := natsgo.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connect to NATS %s: %w", url, err)
	}

	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("create JetStream context: %w", err)
	}

	return &Client{conn: conn, js: js}, nil
}

func (c *Client) Publish(subject string, data []byte) error {
	_, err := c.js.Publish(subject, data)
	return err
}

func (c *Client) JetStream() natsgo.JetStreamContext {
	return c.js
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
