package nats

import (
	"errors"
	"fmt"
	"log"

	"github.com/cashback-platform/services/cashback-service-api/internal/config"
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

	// Create streams if they don't exist
	if err := createStreams(js); err != nil {
		conn.Close()
		return nil, err
	}

	log.Println("NATS connected successfully")
	return &NATSClient{
		conn: conn,
		js:   js,
	}, nil
}

func createStreams(js nats.JetStreamContext) error {
	streams := []struct {
		name     string
		subjects []string
	}{
		{
			name:     "PURCHASE_EVENTS",
			subjects: []string{"purchase.>"},
		},
		{
			name:     "CASHBACK_EVENTS",
			subjects: []string{"cashback.>"},
		},
		{
			name:     "TOKEN_EVENTS",
			subjects: []string{"token.>"},
		},
	}

	for _, s := range streams {
		_, err := js.StreamInfo(s.name)
		if errors.Is(err, nats.ErrStreamNotFound) {
			_, err = js.AddStream(&nats.StreamConfig{
				Name:      s.name,
				Subjects:  s.subjects,
				Retention: nats.LimitsPolicy,
				MaxAge:    7 * 24 * 60 * 60 * 1000000000, // 7 days in nanoseconds
				Storage:   nats.FileStorage,
				Replicas:  1,
			})
			if err != nil {
				return fmt.Errorf("failed to create stream %s: %w", s.name, err)
			}
			log.Printf("Stream %s created", s.name)
		} else if err != nil {
			return fmt.Errorf("failed to get stream info for %s: %w", s.name, err)
		}
	}

	return nil
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
