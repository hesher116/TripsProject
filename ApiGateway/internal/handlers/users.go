package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hesher116/MyFinalProject/ApiGateway/pkg/models" // Впевніться, що шлях правильний
	"github.com/nats-io/nats.go"
)

type UserHandler struct {
	NatsClient *nats.Conn
}

// NewUserHandler створює новий UserHandler з підключенням до NATS
func NewUserHandler(nc *nats.Conn) *UserHandler {
	return &UserHandler{NatsClient: nc}
}

// UserCreate обробляє запит на створення користувача
func (uh *UserHandler) UserCreate(c *gin.Context) {
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

	msg, err := uh.NatsClient.Request("user.create", data, 10*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UserEdit обробляє запит на редагування користувача
func (uh *UserHandler) UserEdit(c *gin.Context) {
	var updateData struct {
		ID          string `json:"id"`
		Username    string `json:"username,omitempty"`
		Email       string `json:"email,omitempty"`
		OldPassword string `json:"oldPassword,omitempty"`
		NewPassword string `json:"newPassword,omitempty"`
	}
	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	data, err := json.Marshal(updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}

	msg, err := uh.NatsClient.Request("user.edit", data, 10*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit user"})
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal response"})
		return
	}

	c.JSON(http.StatusOK, response)
}
