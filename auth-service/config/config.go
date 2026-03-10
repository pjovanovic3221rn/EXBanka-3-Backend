package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort                     string
	DBHost                      string
	DBPort                      string
	DBUser                      string
	DBPassword                  string
	DBName                      string
	DBSSLMode                   string
	JWTSecret                   string
	JWTAccessExpirationMinutes  string
	JWTRefreshExpirationHours   string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("warning: .env file not found, using system environment variables")
	}

	cfg := &Config{
	AppPort:                    getEnv("APP_PORT", "8080"),
	DBHost:                     getEnv("DB_HOST", "localhost"),
	DBPort:                     getEnv("DB_PORT", "5432"),
	DBUser:                     getEnv("DB_USER", "bank"),
	DBPassword:                 getEnv("DB_PASSWORD", "bank123"),
	DBName:                     getEnv("DB_NAME", "bankdb"),
	DBSSLMode:                  getEnv("DB_SSLMODE", "disable"),
	JWTSecret:                  getEnv("JWT_SECRET", "super_secret_key_change_this"),
	JWTAccessExpirationMinutes: getEnv("JWT_ACCESS_EXPIRATION_MINUTES", "15"),
	JWTRefreshExpirationHours:  getEnv("JWT_REFRESH_EXPIRATION_HOURS", "168"),
}

	return cfg
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}