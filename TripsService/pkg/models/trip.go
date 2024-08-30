package models

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Trip struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DriverID    primitive.ObjectID `json:"driverId" bson:"driver_id"`
	PassengerID primitive.ObjectID `json:"passengerId" bson:"passenger_id"`
	Start       string             `json:"start" bson:"start"`
	End         string             `json:"end" bson:"end"`
	Status      string             `json:"status" bson:"status"`
}

// ValidateTrip
func ValidateTrip(trip *Trip) error {
	if trip.DriverID.IsZero() {
		return errors.New("driver ID is required")
	}
	if trip.PassengerID.IsZero() {
		return errors.New("passenger ID is required")
	}
	if trip.Start == "" {
		return errors.New("start location is required")
	}
	if trip.End == "" {
		return errors.New("end location is required")
	}
	if trip.Status == "" {
		return errors.New("status is required")
	}
}
