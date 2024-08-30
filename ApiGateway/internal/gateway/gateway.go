package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hesher116/MyFinalProject/ApiGateway/internal/broker/nats/subjects"
	"github.com/nats-io/nats.go"
)

type GatewayModule struct {
	nats *nats.Conn
}

// NewGatewayModule створює новий GatewayModule з підключенням до NATS
func NewGatewayModule(natsCLI *nats.Conn) *GatewayModule {
	return &GatewayModule{
		nats: natsCLI,
	}
}

// InitNatsSubscribers ініціалізує підписників NATS для обробки подій
func (gm *GatewayModule) InitNatsSubscribers() (err error) {
	_, err = gm.nats.Subscribe(subjects.UserRegEvent.ToString(), gm.HandleRegisterNats)
	if err != nil {
		return fmt.Errorf("failed to subscribe to UserRegEvent: %w", err)
	}

	_, err = gm.nats.Subscribe(subjects.UserAuthEvent.ToString(), gm.HandleAuthorizationNats)
	if err != nil {
		return fmt.Errorf("failed to subscribe to UserAuthEvent: %w", err)
	}

	return nil
}

// RegisterUserNats обробляє реєстрацію користувача через HTTP-запит
func (gm *GatewayModule) RegisterUserNats(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON trouble: " + err.Error()})
		return
	}

	response, err := gm.nats.Request("UserCreateEvent", encode(jsonData), nats.DefaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration process trouble: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, decode(response.Data))
}

// AuthUserNats обробляє авторизацію користувача через HTTP-запит
func (gm *GatewayModule) AuthUserNats(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON(AuthUserNats) trouble: " + err.Error()})
		return
	}

	response, err := gm.nats.Request("UserAuthEvent", encode(jsonData), nats.DefaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authorization process trouble: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, decode(response.Data))
}

// CreateTripNats обробляє створення поїздки через HTTP-запит
func (gm *GatewayModule) CreateTripNats(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON trouble: " + err.Error()})
		return
	}

	response, err := gm.nats.Request("TripCreateEvent", encode(jsonData), nats.DefaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "trip creation process trouble: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, decode(response.Data))
}

// UpdateTripNats обробляє оновлення поїздки через HTTP-запит
func (gm *GatewayModule) UpdateTripNats(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON trouble: " + err.Error()})
		return
	}

	response, err := gm.nats.Request("TripUpdateEvent", encode(jsonData), nats.DefaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "trip update process trouble: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, decode(response.Data))
}

// DeleteTripNats обробляє видалення поїздки через HTTP-запит
func (gm *GatewayModule) DeleteTripNats(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON trouble: " + err.Error()})
		return
	}

	response, err := gm.nats.Request("TripDeleteEvent", encode(jsonData), nats.DefaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "trip deletion process trouble: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, decode(response.Data))
}

// GetTripNats обробляє отримання інформації про поїздку через HTTP-запит
func (gm *GatewayModule) GetTripNats(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON trouble: " + err.Error()})
		return
	}

	response, err := gm.nats.Request("TripGetEvent", encode(jsonData), nats.DefaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "trip retrieval process trouble: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, decode(response.Data))
}

// HandleRegisterNats обробляє реєстрацію користувача через NATS повідомлення
func (gm *GatewayModule) HandleRegisterNats(msg *nats.Msg) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(msg.Data, &jsonData); err != nil {
		fmt.Println("JSON trouble:", err)
		return
	}

	// Обробка отриманих даних та формування відповіді
	response, err := gm.nats.Request("UserService.Register", encode(jsonData), nats.DefaultTimeout)
	if err != nil {
		fmt.Println("registration process trouble:", err)
		return
	}

	_ = msg.Respond(response.Data)
}

// HandleAuthorizationNats обробляє авторизацію користувача через NATS повідомлення
func (gm *GatewayModule) HandleAuthorizationNats(msg *nats.Msg) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(msg.Data, &jsonData); err != nil {
		fmt.Println("JSON trouble:", err)
		return
	}

	// Обробка отриманих даних та формування відповіді
	response, err := gm.nats.Request("UserService.Authorize", encode(jsonData), nats.DefaultTimeout)
	if err != nil {
		fmt.Println("authorization process trouble:", err)
		return
	}

	_ = msg.Respond(response.Data)
}

// encode перетворює структуру в байти для передачі через NATS
func encode(data interface{}) []byte {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error transform to bytes", err)
		return nil
	}
	return bytes
}

// decode перетворює байти з NATS назад у структуру
func decode(data []byte) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("Error transform from bytes", err)
		return nil
	}
	return result
}
