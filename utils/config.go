package utils

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the configuration values for the application
type Config struct {
	Port               string
	DBHost             string
	DBPort             int
	DBUser             string
	DBPassword         string
	DBName             string
	JWTSecret          string
	JWTExpirationHours int
}

// LoadConfig loads environment variables into a Config struct
// It returns an error if loading the .env file fails
func LoadConfig() (*Config, error) {
	// Load .env file if it exists, ignore if it doesn't
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values or environment variables.")
	}

	config := &Config{
		Port:               ":" + getEnv("PORT", "8080"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnvAsInt("DB_PORT", 5432),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "goblog"),
		DBName:             getEnv("DB_NAME", "goblog"),
		JWTSecret:          getEnv("JWT_SECRET", "mysecret"),
		JWTExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 72),
	}

	return config, nil
}

// getEnv retrieves the value of the environment variable named by the key.
// If the variable is not present, it returns the defaultValue.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves the value of the environment variable named by the key as an integer.
// If the variable is not present or cannot be converted to an integer, it returns the defaultValue.
func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
