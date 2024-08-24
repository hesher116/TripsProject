package authorization

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hesher116/MyFinalProject/AuthServsce/internal/broker/nats/subjects"
	"github.com/hesher116/MyFinalProject/AuthServsce/pkg/models"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type AuthorizationModule struct {
	db    *mongo.Client
	cache *redis.Client
	nats  *nats.Conn
}

func NewAuthorizationModule(dbCLI *mongo.Client, cacheCLI *redis.Client, natsCLI *nats.Conn) *AuthorizationModule {
	return &AuthorizationModule{
		db:    dbCLI,
		cache: cacheCLI,
		nats:  natsCLI,
	}
}

func (am *AuthorizationModule) InitNatsSubscribers() (err error) {
	_, err = am.nats.Subscribe(subjects.UserRegister.ToString(), am.RegisterNats)
	if err != nil {
		return err
	}

	_, err = am.nats.Subscribe(subjects.UserAuthorization.ToString(), am.AuthorizationNats)
	if err != nil {
		return err
	}

	return
}

func (am *AuthorizationModule) RegisterNats(m *nats.Msg) {
	log.Print("subscribe succsessfull")
	var user models.User
	err := json.Unmarshal(m.Data, &user)
	if err != nil {
		err = am.nats.Publish(m.Reply, []byte(fmt.Sprintf(`{"error": "Invalid data: %v"}`, err)))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	port := os.Getenv("MONGO_PORT")
	host := os.Getenv("MONGO_HOST")
	mongoUrl := fmt.Sprintf("mongodb://%s:%s", host, port)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		log.Printf("Error connecting to MongoDB: %v", err)
	}

	DB := client.Database("project")
	log.Printf("Connected to MongoDB")

	// Перевірка, чи користувач вже існує в кеші
	cachedUser, err := am.get(ctx, user.ID)
	if err == nil && cachedUser != "" {
		am.nats.Publish(m.Reply, []byte(`{"error": "User already exists in cache"}`))
		return
	}

	// Вставка даних користувача в MongoDB
	_, err = DB.Collection("users").InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			am.nats.Publish(m.Reply, []byte(`{"error": "User already exists"}`))
			return
		}
		am.nats.Publish(m.Reply, []byte(fmt.Sprintf(`{"error": "Failed to register user: %v"}`, err)))
		return
	}

	// Кешування даних користувача в Redis
	userData, _ := json.Marshal(user)
	err = am.set(ctx, user.ID, userData)
	if err != nil {
		log.Printf("Failed to cache user data: %v", err)
	}

	responseData, _ := json.Marshal(map[string]string{
		"status": "User registered successfully",
		"userID": user.ID,
	})

	am.nats.Publish(m.Reply, responseData)
}

func (am *AuthorizationModule) AuthorizationNats(m *nats.Msg) {
	fmt.Printf("AuthorizationNATS called: %s\n", string(m.Data))

	var requestData map[string]string
	err := json.Unmarshal(m.Data, &requestData)
	if err != nil {
		_ = m.Respond([]byte("Invalid request data"))
		return
	}

	userID := requestData["userID"]
	if userID == "" {
		_ = m.Respond([]byte("UserID is missing"))
		return
	}

	ctx := context.Background()

	// 1. Спроба отримати дані користувача з Redis
	userRedisData, err := am.get(ctx, userID)
	if err == nil && userRedisData != "" {
		_ = m.Respond([]byte(userRedisData))
		return
	}

	// 2. Якщо в Redis даних немає, спробуємо знайти їх в MongoDB
	var userMongoData bson.M
	err = am.db.Database("project").Collection("users").FindOne(ctx, bson.D{{"_id", userID}}).Decode(&userMongoData)
	if err == mongo.ErrNoDocuments {
		_ = m.Respond([]byte("User not found in Database"))
		return
	} else if err != nil {
		log.Printf("Error querying MongoDB: %v", err)
		_ = m.Respond([]byte(fmt.Sprintf("Error querying Database: %v", err)))
		return
	}

	// 3. Якщо користувача знайдено в MongoDB, збережемо його в Redis
	userData, err := json.Marshal(userMongoData)
	if err != nil {
		log.Printf("Error marshalling user data: %v", err)
		_ = m.Respond([]byte("Internal server error"))
		return
	}

	err = am.set(ctx, userID, userData, 10*time.Minute) // Зберігаємо в кеші на 10 хвилин
	if err != nil {
		log.Printf("Error caching user data: %v", err)
	}

	// 4. Відправляємо відповідь з даними користувача
	_ = m.Respond(userData)
}

func (am *AuthorizationModule) set(ctx context.Context, key string, value any, expiration ...time.Duration) error {
	var exp time.Duration
	if len(expiration) > 0 {
		exp = expiration[0]
	}
	return am.cache.Set(ctx, key, value, exp).Err()
}

func (am *AuthorizationModule) get(ctx context.Context, key string) (string, error) {
	return am.cache.Get(ctx, key).Result()
}

func (am *AuthorizationModule) del(ctx context.Context, key string) error {
	return am.cache.Del(ctx, key).Err()
}
