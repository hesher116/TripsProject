package nats

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func NewNatsClient(url string) (*nats.Conn, error) {
	nc, err := nats.Connect(url)
	if err != nil {

		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	log.Println("Connected to NATS!")
	return nc, nil
}
