package main

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/cashback-platform/kit/events"
	natsgo "github.com/nats-io/nats.go"
)

type stream struct {
	name    string
	subject string
}

var streams = []stream{
	{
		name:    events.StreamPurchaseEvents,
		subject: events.SubjectPurchaseAll,
	},
	{
		name:    events.StreamCashbackEvents,
		subject: events.SubjectCashbackAll,
	},
	{
		name:    events.StreamTokenEvents,
		subject: events.SubjectTokenAll,
	},
}

func provision(js natsgo.JetStreamContext) error {
	for _, s := range streams {
		if err := ensureStream(js, s); err != nil {
			return fmt.Errorf("setup stream %s: %w", s.name, err)
		}
	}
	return nil
}

func ensureStream(js natsgo.JetStreamContext, s stream) error {
	_, err := js.StreamInfo(s.name)
	if err == nil {
		slog.Info("stream already exists", "stream", s.name)
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

	slog.Info("stream created", "stream", s.name)
	return nil
}
