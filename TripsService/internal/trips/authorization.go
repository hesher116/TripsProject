package trips

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hesher116/MyFinalProject/TripsServsce/internal/broker/nats/subjects"
	"github.com/hesher116/MyFinalProject/TripsServsce/pkg/models"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type TripsModule struct {
	db    *mongo.Client
	nats  *nats.Conn
	redis *redis.Client
}

func NewTripsModule(mongodbCLI *mongo.Client, natsCLI *nats.Conn, redisCLI *redis.Client) *TripsModule {
	return &TripsModule{
		db:    mongodbCLI,
		nats:  natsCLI,
		redis: redisCLI,
	}
}

func (tm *TripsModule) InitNatsSubscribers() (err error) {
	_, err = tm.nats.Subscribe(subjects.TripCreateEvent.ToString(), tm.TripCreateNats)
	if err != nil {
		return err
	}

	_, err = tm.nats.Subscribe(subjects.TripUpdateEvent.ToString(), tm.TripUpdateNats)
	if err != nil {
		return err
	}

	_, err = tm.nats.Subscribe(subjects.TripDeleteEvent.ToString(), tm.TripDeleteNats)
	if err != nil {
		return err
	}

	_, err = tm.nats.Subscribe(subjects.TripGetEvent.ToString(), tm.TripGetNats)
	if err != nil {
		return err
	}

	//_, err = tm.nats.Subscribe(subjects.TripJoinEvent.ToString(), tm.TripJoinNats)
	//if err != nil {
	//	return err
	//}
	//
	//_, err = tm.nats.Subscribe(subjects.TripCancelEvent.ToString(), tm.TripCancelNats)
	//if err != nil {
	//	return err
	//}

	return
}

func (tm *TripsModule) TripCreateNats(m *nats.Msg) {
	var trip models.Trip
	err := json.Unmarshal(m.Data, &trip)
	if err != nil {
		tm.respondWithError(m, fmt.Sprintf("Invalid data: %v", err))
		return
	}

	// Валідація даних подорожі
	if err := models.ValidateTrip(&trip); err != nil {
		tm.respondWithError(m, fmt.Sprintf("Validation error: %v", err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Вставка подорожі в MongoDB
	_, err = tm.db.Database("project").Collection("trips").InsertOne(ctx, trip)
	if err != nil {
		tm.respondWithError(m, fmt.Sprintf("Failed to insert trip: %v", err))
		return
	}

	// Кешування подорожі в Redis
	cacheKey := fmt.Sprintf("trip:%s", trip.ID.Hex())
	tripJson, _ := json.Marshal(trip)
	tm.redis.Set(ctx, cacheKey, tripJson, 10*time.Minute)

	// Відповідь про успішне створення подорожі
	response, _ := json.Marshal(map[string]string{"status": "Trip created successfully", "tripID": trip.ID.Hex()})
	tm.nats.Publish(m.Reply, response)
}

func (tm *TripsModule) TripUpdateNats(m *nats.Msg) {
	var trip models.Trip
	err := json.Unmarshal(m.Data, &trip)
	if err != nil {
		tm.respondWithError(m, fmt.Sprintf("Invalid data: %v", err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": trip.ID}
	update := bson.M{"$set": trip}

	_, err = tm.db.Database("project").Collection("trips").UpdateOne(ctx, filter, update)
	if err != nil {
		tm.respondWithError(m, fmt.Sprintf("Failed to update trip: %v", err))
		return
	}

	// Оновлення кешу в Redis
	cacheKey := fmt.Sprintf("trip:%s", trip.ID.Hex())
	tripJson, _ := json.Marshal(trip)
	tm.redis.Set(ctx, cacheKey, tripJson, 10*time.Minute)

	tm.nats.Publish(m.Reply, []byte("Trip updated successfully"))
}

func (tm *TripsModule) TripGetNats(m *nats.Msg) {
	var tripID primitive.ObjectID
	err := json.Unmarshal(m.Data, &tripID)
	if err != nil {
		tm.respondWithError(m, fmt.Sprintf("Invalid data: %v", err))
		return
	}

	ctx := context.Background()

	// Перевірка в Redis
	cacheKey := fmt.Sprintf("trip:%s", tripID.Hex())
	cachedTrip, err := tm.redis.Get(ctx, cacheKey).Result()
	if err == nil && cachedTrip != "" {
		_ = m.Respond([]byte(cachedTrip))
		return
	}

	// Якщо в кеші немає, шукаємо в MongoDB
	var trip models.Trip
	err = tm.db.Database("project").Collection("trips").FindOne(ctx, bson.M{"_id": tripID}).Decode(&trip)
	if err != nil {
		tm.respondWithError(m, fmt.Sprintf("Trip not found: %v", err))
		return
	}

	response, _ := json.Marshal(trip)

	// Збереження в кеш Redis
	tm.redis.Set(ctx, cacheKey, response, 10*time.Minute)

	tm.nats.Publish(m.Reply, response)
}

func (tm *TripsModule) TripDeleteNats(m *nats.Msg) {
	var tripID primitive.ObjectID
	err := json.Unmarshal(m.Data, &tripID)
	if err != nil {
		tm.respondWithError(m, fmt.Sprintf("Invalid data: %v", err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": tripID}

	_, err = tm.db.Database("project").Collection("trips").DeleteOne(ctx, filter)
	if err != nil {
		tm.respondWithError(m, fmt.Sprintf("Failed to delete trip: %v", err))
		return
	}

	// Видалення подорожі з кешу Redis
	cacheKey := fmt.Sprintf("trip:%s", tripID.Hex())
	tm.redis.Del(ctx, cacheKey)

	tm.nats.Publish(m.Reply, []byte("Trip deleted successfully"))
}

func (tm *TripsModule) respondWithError(m *nats.Msg, errorMsg string) {
	log.Println(errorMsg)
	_ = m.Respond([]byte(fmt.Sprintf(`{"error": "%s"}`, errorMsg)))
}
