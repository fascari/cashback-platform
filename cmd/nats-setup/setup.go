package main

import (
	"fmt"
	"log/slog"
	"os"

	natsgo "github.com/nats-io/nats.go"
)

type config struct {
	url string
}

func loadConfig() config {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = natsgo.DefaultURL
	}
	return config{url: url}
}

func run() error {
	cfg := loadConfig()

	conn, err := natsgo.Connect(cfg.url)
	if err != nil {
		return fmt.Errorf("connect to NATS %s: %w", cfg.url, err)
	}
	defer conn.Close()

	js, err := conn.JetStream()
	if err != nil {
		return fmt.Errorf("create JetStream context: %w", err)
	}

	if err := provision(js); err != nil {
		return err
	}

	slog.Info("all streams ready")
	return nil
}
