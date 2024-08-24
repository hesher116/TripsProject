package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	natsclient "github.com/hesher116/MyFinalProject/ApiGateway/internal/broker/nats"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"net/http"
)

func main() {
	godotenv.Load()

	nc, err := natsclient.NewNatsClient()
	if err != nil {
		fmt.Println("Error Connected to NATS:", err)
		return
	}
	defer nc.Close()

	fmt.Println("Connected to NATS")

	// Ініціалізація HTTP-сервера з використанням Gin
	router := gin.Default()

	// Маршрутизація запитів
	router.POST("/register", func(c *gin.Context) {
		RegisterUserNats(c, nc)
		fmt.Println("Register User Nats Success")
	})

	router.POST("/login", func(c *gin.Context) {
		AuthUserNats(c, nc)
		fmt.Println("Login User Nats Success")
	})

	// Запуск HTTP-сервера
	router.Run(":8080")

	select {}
}

// RegisterUserNats обробляє реєстрацію користувача
func RegisterUserNats(c *gin.Context, nc *nats.Conn) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := nc.Request("UserCreateEvent", encode(jsonData), nats.DefaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, decode(response.Data))
}

// AuthUserNats обробляє авторизацію користувача
func AuthUserNats(c *gin.Context, nc *nats.Conn) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := nc.Request("UserAuthEvent", encode(jsonData), nats.DefaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, decode(response.Data))
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
