package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	VehicleDBDSN string
	UserDBDSN     string
	AppEnv        string
	JWTSigningKey string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	return &Config{
		Port:         getEnv("PORT", "8080"),
		VehicleDBDSN:  getEnv("VEHICLE_DB_DSN", ""),
		UserDBDSN:     getEnv("FMS_USER_DB_DSN", ""),
		AppEnv:        getEnv("APP_ENV", "development"),
		JWTSigningKey: getEnv("JWT_SIGNING_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
