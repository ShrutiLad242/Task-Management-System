package config

import (
	"fmt"
	"os"
	"time"
	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	JWTSecret         string
	AutoCompleteDelay time.Duration
}

func Load() *Config {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file not found, using system env")
	}

	delay, err := time.ParseDuration(getEnv("AUTO_COMPLETE_DELAY", "1m"))
	if err != nil {
		fmt.Print("Invalid AUTO_COMPLETE_DELAY")
	}

	return &Config{
		Port:              getEnv("PORT", "8080"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "taskdb"),
		JWTSecret:         getEnv("JWT_SECRET", "secret"),
		AutoCompleteDelay: delay,
	}
}

func getEnv(key string, fallback string) string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
