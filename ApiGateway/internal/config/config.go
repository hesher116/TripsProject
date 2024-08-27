package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	NatsURL string
	Envi    string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Error loading .env file", err)
	}

	config := &Config{
		NatsURL: getEnv("NATS_URL", "nats://nats:4222"),
		Envi:    getEnv("ENVIRONMENT", "prod"),
	}

	return config
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
