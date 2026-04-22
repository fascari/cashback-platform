package main

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/cashback-platform/kit/events"
	natsgo "github.com/nats-io/nats.go"
)

type consumer struct {
	stream string
	config *natsgo.ConsumerConfig
}

var consumers = []consumer{
	{
		stream: events.StreamDepositEvents,
		config: &natsgo.ConsumerConfig{
			Durable:       "cashback-service-api-deposits",
			FilterSubject: events.DepositDetected,
			AckPolicy:     natsgo.AckExplicitPolicy,
			MaxDeliver:    5,
			AckWait:       30 * time.Second,
			DeliverPolicy: natsgo.DeliverAllPolicy,
		},
	},
}

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
	{
		name:    events.StreamDepositEvents,
		subject: events.SubjectDepositAll,
	},
}

func provision(js natsgo.JetStreamContext) error {
	for _, s := range streams {
		if err := ensureStream(js, s); err != nil {
			return fmt.Errorf("setup stream %s: %w", s.name, err)
		}
	}
	for _, c := range consumers {
		if err := ensureConsumer(js, c); err != nil {
			return fmt.Errorf("setup consumer %s: %w", c.config.Durable, err)
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

func ensureConsumer(js natsgo.JetStreamContext, c consumer) error {
	_, err := js.ConsumerInfo(c.stream, c.config.Durable)
	if err == nil {
		slog.Info("consumer already exists", "consumer", c.config.Durable)
		return nil
	}
	if !errors.Is(err, natsgo.ErrConsumerNotFound) {
		return fmt.Errorf("query consumer info: %w", err)
	}
	if _, err := js.AddConsumer(c.stream, c.config); err != nil {
		return fmt.Errorf("add consumer: %w", err)
	}
	slog.Info("consumer created", "consumer", c.config.Durable)
	return nil
}
