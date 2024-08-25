package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func NewNatsClient(url string) (*nats.Conn, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	log.Println("Connected to NATS!")
	return nc, nil
}
