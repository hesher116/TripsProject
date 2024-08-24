package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
)

func NewNatsClient() (nc *nats.Conn, err error) {
	fmt.Println("Connected to NATS(package nats)")
	return nats.Connect(os.Getenv("NATS_URL"))
}
