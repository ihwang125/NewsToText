package config

import (
	"os"
)

type Config struct {
	Environment   string
	Port          string
	DatabaseURL   string
	RedisURL      string
	JWTSecret     string
	NewsAPIKey    string
	SMSAPIKey     string
	LogLevel      string
}

func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "root:password@tcp(localhost:3306)/newstotext?charset=utf8mb4&parseTime=True&loc=Local"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:   getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		NewsAPIKey:  getEnv("NEWS_API_KEY", ""),
		SMSAPIKey:   getEnv("SMS_API_KEY", ""),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}