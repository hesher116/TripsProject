// internal/storage/mongo/mongo.go
package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func Connect(ctx context.Context, host, port string) (*mongo.Client, error) {
	// Формуємо URI для підключення
	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port)

	// Створюємо опції клієнта на основі URI
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Використовуємо ці опції для підключення до MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}
