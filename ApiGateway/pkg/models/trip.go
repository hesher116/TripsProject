package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Trip представляє модель поїздки
type Trip struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`                    // Унікальний ідентифікатор поїздки
	DriverID    primitive.ObjectID `json:"driver_id,omitempty" bson:"driver_id,omitempty"`       // Ідентифікатор водія
	PassengerID primitive.ObjectID `json:"passenger_id,omitempty" bson:"passenger_id,omitempty"` // Ідентифікатор пасажира (якщо є)
	StartPoint  string             `json:"start_point" bson:"start_point"`                       // Початкова точка поїздки
	EndPoint    string             `json:"end_point" bson:"end_point"`                           // Кінцева точка поїздки
	StartTime   time.Time          `json:"start_time" bson:"start_time"`                         // Час початку поїздки
	EndTime     time.Time          `json:"end_time" bson:"end_time"`                             // Час закінчення поїздки (може бути порожнім)
	Status      string             `json:"status" bson:"status"`                                 // Статус поїздки (наприклад, "scheduled", "completed", "cancelled")
}
