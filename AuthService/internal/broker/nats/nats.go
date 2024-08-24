package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

func NewNatsClient() (nc *nats.Conn, err error) {
	fmt.Println("Підключено до NATS")
	return nats.Connect("nats://localhost:4222")
}
