package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Environment   string
	ServerPort    string
	LogLevel      string
	MessageTTL    time.Duration
	MaxRetries    int
	RetryInterval time.Duration
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Environment:   getEnv("APP_ENV", "development"),
		ServerPort:    getEnv("PORT", "8080"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		MessageTTL:    time.Duration(getIntEnv("MESSAGE_TTL_HOURS", 24)) * time.Hour,
		MaxRetries:    getIntEnv("MAX_RETRIES", 3),
		RetryInterval: time.Duration(getIntEnv("RETRY_INTERVAL_SECONDS", 60)) * time.Second,
	}

	log.Printf("Loaded configuration: Environment=%s, Port=%s, LogLevel=%s, MessageTTL=%v, MaxRetries=%d, RetryInterval=%v",
		config.Environment, config.ServerPort, config.LogLevel, config.MessageTTL, config.MaxRetries, config.RetryInterval)

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
