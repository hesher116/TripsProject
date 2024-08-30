package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hesher116/MyFinalProject/ApiGateway/pkg/models" // Впевніться, що шлях правильний
	nats "github.com/nats-io/nats.go"
)

type TripsHandler struct {
	NatsClient *nats.Conn
}

func NewTripsHandler(nc *nats.Conn) *TripsHandler {
	return &TripsHandler{NatsClient: nc}
}

// CreateTrip обробляє запит на створення поїздки
func (th *TripsHandler) CreateTrip(c *gin.Context) {
	var trip models.Trip
	if err := c.BindJSON(&trip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	data, err := json.Marshal(trip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}

	msg, err := th.NatsClient.Request("trip.create", data, 10*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trip"})
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateTrip обробляє запит на оновлення поїздки
func (th *TripsHandler) UpdateTrip(c *gin.Context) {
	var trip models.Trip
	if err := c.BindJSON(&trip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Конвертація даних у JSON
	data, err := json.Marshal(trip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}

	// Відправка запиту через NATS
	msg, err := th.NatsClient.Request("trip.update", data, 10*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trip"})
		return
	}

	// Обробка відповіді від TripsService
	var response map[string]interface{}
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTrip обробляє запит на отримання поїздки за ID
func (th *TripsHandler) GetTrip(c *gin.Context) {
	var tripID struct {
		ID string `json:"id"`
	}
	if err := c.BindJSON(&tripID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	data, err := json.Marshal(tripID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}

	msg, err := th.NatsClient.Request("trip.get", data, 10*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trip"})
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteTrip обробляє запит на видалення поїздки
func (th *TripsHandler) DeleteTrip(c *gin.Context) {
	var tripID struct {
		ID string `json:"id"`
	}
	if err := c.BindJSON(&tripID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	data, err := json.Marshal(tripID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}

	msg, err := th.NatsClient.Request("trip.delete", data, 10*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete trip"})
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// JoinTrip обробляє запит на приєднання до поїздки
func (th *TripsHandler) JoinTrip(c *gin.Context) {
	var joinData struct {
		TripID string `json:"trip_id"`
		UserID string `json:"user_id"`
	}
	if err := c.BindJSON(&joinData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	data, err := json.Marshal(joinData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}

	msg, err := th.NatsClient.Request("trip.join", data, 10*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join trip"})
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CancelTrip обробляє запит на скасування поїздки
func (th *TripsHandler) CancelTrip(c *gin.Context) {
	var cancelData struct {
		TripID string `json:"trip_id"`
		UserID string `json:"user_id"`
	}
	if err := c.BindJSON(&cancelData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	data, err := json.Marshal(cancelData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}

	msg, err := th.NatsClient.Request("trip.cancel", data, 10*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel trip"})
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response"})
		return
	}

	c.JSON(http.StatusOK, response)
}
