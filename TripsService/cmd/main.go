package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hesher116/MyFinalProject/TripsServsce/internal/trips"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"

	"github.com/hesher116/MyFinalProject/TripsServsce/internal/broker/nats"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Помилка завантаження .env файлу")
	}
	for _, e := range os.Environ() {
		fmt.Println(e)
	}

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Помилка підключення REDIS:", err)
		return
	}
	fmt.Println("Підключено до Redis:", pong)

	port := os.Getenv("MONGO_PORT")
	host := os.Getenv("MONGO_HOST")
	if port == "" || host == "" {
		log.Fatal("Не видно MONGO_HOST або MONGO_PORT")
	}
	mongoUrl := fmt.Sprintf("mongodb://%s:%s", host, port)

	clientOptions := options.Client().ApplyURI(mongoUrl)
	cli, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Disconnect(context.TODO())

	err = cli.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	nc, err := nats.NewNatsClient()
	if err != nil {
		fmt.Println("Помилка підключення NATS:", err)
		return
	}
	defer nc.Close()

	tripsModule := trips.NewTripsModule(cli, nc, rdb)
	tripsModule.InitNatsSubscribers()

	select {}
}
