package main

import (
	"context"
	"fmt"
	"github.com/hesher116/MyFinalProject/UserServsce/internal/users"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"

	"github.com/hesher116/MyFinalProject/UserServsce/internal/broker/nats"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Помилка завантаження .env файлу")
	}
	for _, e := range os.Environ() {
		fmt.Println(e)
	}

	port := os.Getenv("MONGO_PORT")
	host := os.Getenv("MONGO_HOST")
	mongoUrl := fmt.Sprintf("mongodb://%s:%s", host, port)

	clientOptions := options.Client().ApplyURI(mongoUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	err = client.Ping(context.TODO(), nil)
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

	userModule := users.NewUserModule(client, nc)
	userModule.InitNatsSubscribers()

	select {}
}
