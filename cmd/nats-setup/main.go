// cmd/nats-setup provisions JetStream streams required by the cashback platform.
// Run once during infrastructure setup, before starting any service.
// The NATS_URL environment variable controls the target server (default: nats://localhost:4222).
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cashback-platform/kit/events"
	natsgo "github.com/nats-io/nats.go"
)

type streamDef struct {
	name    string
	subject string
}

var streams = []streamDef{
	{events.StreamPurchaseEvents, events.SubjectPurchaseAll},
	{events.StreamCashbackEvents, events.SubjectCashbackAll},
	{events.StreamTokenEvents, events.SubjectTokenAll},
}

func main() {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = natsgo.DefaultURL
	}

	conn, err := natsgo.Connect(url)
	if err != nil {
		log.Fatalf("connect to NATS %s: %v", url, err)
	}
	defer conn.Close()

	js, err := conn.JetStream()
	if err != nil {
		log.Fatalf("create JetStream context: %v", err)
	}

	for _, s := range streams {
		if err := ensureStream(js, s); err != nil {
			log.Fatalf("setup stream %s: %v", s.name, err)
		}
	}

	log.Println("nats-setup: all streams ready")
}

func ensureStream(js natsgo.JetStreamContext, s streamDef) error {
	_, err := js.StreamInfo(s.name)
	if err == nil {
		log.Printf("stream %s already exists", s.name)
		return nil
	}

	if !errors.Is(err, natsgo.ErrStreamNotFound) {
		return fmt.Errorf("query stream info: %w", err)
	}

	_, err = js.AddStream(&natsgo.StreamConfig{
		Name:      s.name,
		Subjects:  []string{s.subject},
		Retention: natsgo.LimitsPolicy,
		MaxAge:    7 * 24 * time.Hour,
		Storage:   natsgo.FileStorage,
		Replicas:  1,
	})
	if err != nil {
		return fmt.Errorf("add stream: %w", err)
	}

	log.Printf("stream %s created", s.name)
	return nil
}
