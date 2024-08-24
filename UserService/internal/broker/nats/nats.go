package nats

import (
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

func NewNatsClient() (nc *nats.Conn, err error) {
	fmt.Println("Підключено до NATS")
	return nats.Connect(os.Getenv("NATS_URL"))
}
