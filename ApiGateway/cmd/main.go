package main

import (
	"fmt"
	"github.com/hesher116/MyFinalProject/ApiGateway/internal/gateway"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	natsclient "github.com/hesher116/MyFinalProject/ApiGateway/internal/broker/nats"
	"github.com/hesher116/MyFinalProject/ApiGateway/internal/config"
	"github.com/nats-io/nats.go"
)

func main() {
	moveToRelease, err := strconv.ParseBool(os.Getenv("MOVE_TO_RELEASE"))
	if err != nil {
		log.Printf("Invalid value for MOVE_TO_RELEASE: %v.", err)
	}

	if moveToRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	cfg := config.LoadConfig()

	natsClient, err := natsclient.NewNatsClient(cfg.NatsURL)
	if err != nil {
		log.Fatalf("Nats error: %v", err)
	}
	defer natsClient.Close()

	// Ініціалізація HTTP-сервера з використанням Gin
	router := gin.Default()

	authModule := gateway.NewAuthorizationModule(natsClient)
	err = authModule.InitNatsSubscribers()
	if err != nil {
		log.Fatalf("Failed to initialize NATS subscribers: %v", err)
	}

	// Маршрутизація запитів
	router.POST("/register", func(c *gin.Context) {
		RegisterUserNats(c, natsClient)
		fmt.Println("Register User Nats Success")
	})

	router.POST("/login", func(c *gin.Context) {
		AuthUserNats(c, natsClient)
		fmt.Println("Login User Nats Success")
	})

	// Запуск HTTP-сервера
	log.Println("Starting server...")
	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	log.Println("Server is running...")

	select {}
}

// RegisterUserNats обробляє реєстрацію користувача
func RegisterUserNats(c *gin.Context, nc *nats.Conn) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON trouble" + err.Error()})
		return
	}

	response, err := nc.Request("UserCreateEvent", encode(jsonData), nats.DefaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration process trouble" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, decode(response.Data))
}

// AuthUserNats обробляє авторизацію користувача
func AuthUserNats(c *gin.Context, nc *nats.Conn) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON(AuthUserNats) trouble" + err.Error()})
		return
	}

	response, err := nc.Request("UserAuthEvent", encode(jsonData), nats.DefaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration process trouble" + err.Error()})
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

//{
//"username":"Maksim",
//"password":"123456"
//}
