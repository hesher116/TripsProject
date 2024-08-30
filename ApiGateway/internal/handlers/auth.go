package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hesher116/MyFinalProject/ApiGateway/pkg/models" // Впевніться, що шлях правильний
	"github.com/nats-io/nats.go"
)

type AuthHandler struct {
	NatsClient *nats.Conn
}

// NewAuthHandler створює новий AuthHandler з підключенням до NATS
func NewAuthHandler(nc *nats.Conn) *AuthHandler {
	return &AuthHandler{NatsClient: nc}
}

// UserRegister обробляє запит на реєстрацію користувача
func (ah *AuthHandler) UserRegister(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}

	msg, err := ah.NatsClient.Request("user.register", data, 10*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UserAuthorization обробляє запит на авторизацію користувача
func (ah *AuthHandler) UserAuthorization(c *gin.Context) {
	var authData map[string]string
	if err := c.BindJSON(&authData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	data, err := json.Marshal(authData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}

	msg, err := ah.NatsClient.Request("user.authorization", data, 10*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authorize user"})
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response"})
		return
	}

	c.JSON(http.StatusOK, response)
}
