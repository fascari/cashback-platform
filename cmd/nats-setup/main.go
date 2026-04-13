package main

import (
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("nats-setup: %v", err)
	}
}
