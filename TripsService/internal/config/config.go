package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	MongoHost string
	MongoPort string
	RedisURL  string
	NatsURL   string
	Envi      string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &Config{
		MongoHost: getEnv("MONGO_HOST", "localhost"),
		MongoPort: getEnv("MONGO_PORT", "27017"),
		RedisURL:  getEnv("REDIS_URL", "localhost:6379"),
		NatsURL:   getEnv("NATS_URL", "nats://localhost:4222"),
		Envi:      getEnv("ENVIRONMENT", "prod"),
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
