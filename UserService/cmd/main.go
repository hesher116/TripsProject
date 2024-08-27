package main

import (
	"context"
	"github.com/hesher116/MyFinalProject/UserService/internal/broker/nats"
	"github.com/hesher116/MyFinalProject/UserService/internal/config"
	"github.com/hesher116/MyFinalProject/UserService/internal/storage/mongo"
	"github.com/hesher116/MyFinalProject/UserService/internal/users"
	"log"
)

func main() {

	cfg := config.LoadConfig()

	ctx := context.Background()

	mongoClient, err := mongo.Connect(ctx, cfg.MongoHost, cfg.MongoPort)
	if err != nil {
		log.Fatalf("MongoDB error: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	natsClient, err := nats.NewNatsClient(cfg.NatsURL)
	if err != nil {
		log.Fatalf("Nats error: %v", err)
	}
	defer natsClient.Close()

	userModule := users.NewUserModule(mongoClient, natsClient)
	err = userModule.InitNatsSubscribers()
	if err != nil {
		log.Fatalf("Failed to initialize NATS subscribers: %v", err)
	}

	log.Println("Initialized NATS subscribers...")

	select {}
}
