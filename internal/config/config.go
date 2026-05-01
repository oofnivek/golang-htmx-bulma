package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	VehicleDBDSN string
	AppEnv       string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	return &Config{
		Port:         getEnv("PORT", "8080"),
		VehicleDBDSN: getEnv("VEHICLE_DB_DSN", ""),
		AppEnv:       getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
