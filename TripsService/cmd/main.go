package main

import (
	"context"
	"log"

	"github.com/hesher116/MyFinalProject/TripsService/internal/broker/nats"
	"github.com/hesher116/MyFinalProject/TripsService/internal/config"
	"github.com/hesher116/MyFinalProject/TripsService/internal/storage/mongo"
	"github.com/hesher116/MyFinalProject/TripsService/internal/storage/redis"
	"github.com/hesher116/MyFinalProject/TripsService/internal/trips"
)

func main() {

	cfg := config.LoadConfig()

	ctx := context.Background()

	mongoClient, err := mongo.Connect(ctx, cfg.MongoHost, cfg.MongoPort)
	if err != nil {
		log.Fatalf("MongoDB error: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	redisClient, err := redis.Connect(ctx, cfg.RedisURL)
	if err != nil {
		log.Fatalf("Redis error: %v", err)
	}
	defer redisClient.Close()

	natsClient, err := nats.NewNatsClient(cfg.NatsURL)
	if err != nil {
		log.Fatalf("Nats error: %v", err)
	}
	defer natsClient.Close()

	tripsModule := trips.NewTripsModule(mongoClient, redisClient, natsClient)
	err = tripsModule.InitNatsSubscribers()
	if err != nil {
		log.Fatalf("Failed to initialize NATS subscribers: %v", err)
	}

	log.Println("Initialized NATS subscribers...")

	select {}
}
