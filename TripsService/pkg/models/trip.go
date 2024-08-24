package models

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Trip struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DriverID    primitive.ObjectID `json:"driverId" bson:"driver_id"`
	PassengerID primitive.ObjectID `json:"passengerId" bson:"passenger_id"`
	Start       string             `json:"start" bson:"start"`
	End         string             `json:"end" bson:"end"`
	Status      string             `json:"status" bson:"status"`
}

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
	if trip.Status != "active" && trip.Status != "completed" && trip.Status != "canceled" {
		return errors.New("status must be one of the following: active, completed, canceled")
	}

	startTime, err := time.Parse(time.RFC3339, trip.Start)
	if err != nil {
		return errors.New("invalid start date format")
	}

	endTime, err := time.Parse(time.RFC3339, trip.End)
	if err != nil {
		return errors.New("invalid end date format")
	}

	if startTime.After(endTime) {
		return errors.New("start date must be before or equal to end date")
	}

	return nil
}
