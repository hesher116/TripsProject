package main

import (
	"github.com/hesher116/MyFinalProject/ApiGateway/internal/gateway"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	natsclient "github.com/hesher116/MyFinalProject/ApiGateway/internal/broker/nats"
	"github.com/hesher116/MyFinalProject/ApiGateway/internal/config"
)

func main() {
	godotenv.Load()
	moveToRelease, err := strconv.ParseBool(os.Getenv("MOVE_TO_RELEASE"))
	if err != nil {
		log.Printf("Invalid value for MOVE_TO_RELEASE: %v.", err)
	}

	if moveToRelease {
		log.Printf("MOVE_TO_RELEASE value: %s", os.Getenv("MOVE_TO_RELEASE"))
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

	gatewayModule := gateway.NewGatewayModule(natsClient)
	err = gatewayModule.InitNatsSubscribers()
	if err != nil {
		log.Fatalf("Failed to initialize NATS subscribers: %v", err)
	}

	//// Маршрутизація запитів
	//router.POST("/register", func(c *gin.Context) {
	//	gatewayModule.RegisterUserNats(c)
	//	fmt.Println("Register User Nats Success")
	//})
	//
	//router.POST("/login", func(c *gin.Context) {
	//	gatewayModule.AuthUserNats(c)
	//	fmt.Println("Login User Nats Success")
	//})

	// Запуск HTTP-сервера
	log.Println("Starting server...")
	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	log.Println("Server is running...")

	select {}
}
