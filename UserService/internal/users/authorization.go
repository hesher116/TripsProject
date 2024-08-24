package users

import (
	"context"
	"encoding/json"
	"github.com/hesher116/MyFinalProject/UserServsce/internal/broker/nats/subjects"
	"github.com/hesher116/MyFinalProject/UserServsce/pkg/models"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type UserModule struct {
	db   *mongo.Client
	nats *nats.Conn
}

func NewUserModule(mongodbCli *mongo.Client, natsCLI *nats.Conn) *UserModule {
	return &UserModule{
		db:   mongodbCli,
		nats: natsCLI,
	}
}

func (um *UserModule) InitNatsSubscribers() (err error) {
	_, err = um.nats.Subscribe(subjects.UserCreateEvent.ToString(), um.UserCreateNats)
	if err != nil {
		return err
	}

	_, err = um.nats.Subscribe(subjects.UserEditEvent.ToString(), um.UserEditNats)
	if err != nil {
		return err
	}

	return
}

func (um *UserModule) UserCreateNats(m *nats.Msg) {
	var user models.User
	err := json.Unmarshal(m.Data, &user)
	if err != nil {
		log.Printf("Error unmarshalling UserCreateEvent: %v", err)
		return
	}

	// Перевірка даних
	if user.Username == "" || user.Email == "" || user.Password == "" {
		log.Println("Invalid user data")
		return
	}

	// Збереження в базу даних
	collection := um.db.Database("project").Collection("users")
	uId, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return
	}
	user.ID = (uId.InsertedID).(primitive.ObjectID)
	log.Printf("User created: %s", user.Username)
}

func (um *UserModule) UserEditNats(m *nats.Msg) {
	var updateData struct {
		ID          primitive.ObjectID `json:"id"`
		Username    string             `json:"username,omitempty"`
		Email       string             `json:"email,omitempty"`
		OldPassword string             `json:"oldPassword,omitempty"`
		NewPassword string             `json:"newPassword,omitempty"`
	}
	err := json.Unmarshal(m.Data, &updateData)
	if err != nil {
		log.Printf("Error unmarshalling UserEditEvent: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := um.db.Database("project").Collection("users")

	// Знайти користувача в базі даних
	var existingUser models.User
	err = collection.FindOne(ctx, bson.M{"_id": updateData.ID}).Decode(&existingUser)
	if err != nil {
		log.Printf("Error finding user in database: %v", err)
		return
	}

	// Перевірка старого пароля
	if updateData.OldPassword != "" && updateData.OldPassword != existingUser.Password {
		log.Println("Old password does not match")
		return
	}

	// Оновлення даних користувача
	update := bson.M{}
	if updateData.Username != "" {
		update["username"] = updateData.Username
	}
	if updateData.Email != "" {
		update["email"] = updateData.Email
	}
	if updateData.NewPassword != "" {
		update["password"] = updateData.NewPassword
	}

	if len(update) > 0 {
		_, err = collection.UpdateOne(ctx, bson.M{"_id": updateData.ID}, bson.M{"$set": update})
		if err != nil {
			log.Printf("Error updating user in database: %v", err)
			return
		}

		log.Printf("User updated: %s", updateData.ID.Hex())
	}
}
