package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	// database configuration
	DB_DSN string

	// JWT configuration
	JWT_SECRET     string
	JWT_EXPIRES_IN string

	// server configuration
	SERVER_HOST string
	SERVER_PORT string

	// environment configuration
	ENVIRONMENT string

	// log configuration
	LOG_LEVEL string

	// CORS configuration
	CORS_ALLOWED_ORIGINS string
	CORS_ALLOWED_METHODS string
	CORS_ALLOWED_HEADERS string

	// external API configuration
	EXTERNAL_API_BASE_URL string
	EXTERNAL_API_KEY      string

	// password encryption configuration
	BCRYPT_COST int

	//DEEPSEEK API
	OPENAI_API_KEY string
}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		bcryptCost := 12 
		if costStr := os.Getenv("BCRYPT_COST"); costStr != "" {
			if cost, err := strconv.Atoi(costStr); err == nil {
				bcryptCost = cost
			}
		}

		cfg = &Config{
			DB_DSN:                os.Getenv("DB_DSN"),
			JWT_SECRET:            os.Getenv("JWT_SECRET"),
			JWT_EXPIRES_IN:        getEnvWithDefault("JWT_EXPIRES_IN", "24h"),
			SERVER_HOST:           getEnvWithDefault("SERVER_HOST", "localhost"),
			SERVER_PORT:           getEnvWithDefault("SERVER_PORT", "8080"),
			ENVIRONMENT:           getEnvWithDefault("ENVIRONMENT", "development"),
			LOG_LEVEL:             getEnvWithDefault("LOG_LEVEL", "info"),
			CORS_ALLOWED_ORIGINS:  getEnvWithDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
			CORS_ALLOWED_METHODS:  getEnvWithDefault("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
			CORS_ALLOWED_HEADERS:  getEnvWithDefault("CORS_ALLOWED_HEADERS", "Content-Type,Authorization"),
			EXTERNAL_API_BASE_URL: os.Getenv("EXTERNAL_API_BASE_URL"),
			EXTERNAL_API_KEY:      os.Getenv("EXTERNAL_API_KEY"),
			BCRYPT_COST:           bcryptCost,
			OPENAI_API_KEY:        os.Getenv("OPENAI_API_KEY"),
		}

		if cfg.JWT_SECRET == "" {
			log.Fatal("JWT_SECRET is not set in the environment variables")
		}
		if cfg.OPENAI_API_KEY == "" {
			log.Fatal("OPENAI_API_KEY is not set in the environment variables")
		}
	})
}

func GetConfig() *Config {
	return cfg
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
