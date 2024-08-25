package main

import (
	"context"
	"github.com/hesher116/MyFinalProject/AuthServsce/internal/authorization"
	"github.com/hesher116/MyFinalProject/AuthServsce/internal/broker/nats"
	"github.com/hesher116/MyFinalProject/AuthServsce/internal/config"
	"github.com/hesher116/MyFinalProject/AuthServsce/internal/storage/mongo"
	"github.com/hesher116/MyFinalProject/AuthServsce/internal/storage/redis"

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

	authModule := authorization.NewAuthorizationModule(mongoClient, redisClient, natsClient)
	authModule.InitNatsSubscribers()

	log.Println("Server is running...")

	select {}
}
