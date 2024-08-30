package main

import (
	"github.com/hesher116/MyFinalProject/ApiGateway/internal/gateway"
	"github.com/hesher116/MyFinalProject/ApiGateway/internal/handlers"
	"github.com/joho/godotenv"
	"log"
	"net/http"
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
	authHandler := handlers.NewAuthHandler(natsClient)
	tripsHandler := handlers.NewTripsHandler(natsClient)
	userHandler := handlers.NewUserHandler(natsClient)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	http.ListenAndServe(":3000", nil)

	// Визначення маршрутів для Gateway
	router.POST("/gateway/register", gatewayModule.RegisterUserNats)
	router.POST("/gateway/login", gatewayModule.AuthUserNats)
	router.POST("/gateway/trips/create", gatewayModule.CreateTripNats)
	router.POST("/gateway/trips/update", gatewayModule.UpdateTripNats)
	router.POST("/gateway/trips/delete", gatewayModule.DeleteTripNats)
	router.POST("/gateway/trips/get", gatewayModule.GetTripNats)

	// Визначення маршрутів для Authorization
	router.POST("/auth/register", authHandler.UserRegister)
	router.POST("/auth/login", authHandler.UserAuthorization)

	// Визначення маршрутів для Trips
	router.POST("/trips/create", tripsHandler.CreateTrip)
	router.POST("/trips/update", tripsHandler.UpdateTrip)
	router.POST("/trips/get", tripsHandler.GetTrip)
	router.POST("/trips/delete", tripsHandler.DeleteTrip)
	router.POST("/trips/join", tripsHandler.JoinTrip)
	router.POST("/trips/cancel", tripsHandler.CancelTrip)

	// Визначення маршрутів для Users
	router.POST("/users/create", userHandler.UserCreate)
	router.POST("/users/edit", userHandler.UserEdit)

	// Запуск HTTP-сервера
	log.Println("Starting server...")
	err = router.Run(":3000")
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	log.Println("Server is running...")

	select {}
}
