package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hesher116/MyFinalProject/AuthServsce/internal/authorization"
	"github.com/hesher116/MyFinalProject/AuthServsce/internal/broker/nats"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	//err := godotenv.Load("../.env")
	//if err != nil {
	//	log.Fatal("Помилка завантаження .env файлу")
	//}
	//for _, e := range os.Environ() {
	//	fmt.Println(e)
	//}

	//port := os.Getenv("MONGO_PORT")
	//host := os.Getenv("MONGO_HOST")
	//if port == "" || host == "" {
	//	log.Fatal("Не видно MONGO_HOST або MONGO_PORT")
	//}
	//mongoUrl := fmt.Sprintf("mongodb://%s:%s", host, port)

	mongoUrl := fmt.Sprintf("mongodb://mongo:27017")

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

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Помилка підключення REDIS:", err)
		return
	}
	fmt.Println("Підключено до Redis:", pong)

	nc, err := nats.NewNatsClient()
	if err != nil {
		fmt.Println("Помилка підключення NATS:", err)
		return
	}
	defer nc.Close()

	authModule := authorization.NewAuthorizationModule(client, rdb, nc)
	authModule.InitNatsSubscribers()

	log.Println("Server is running...")

	select {}
}
