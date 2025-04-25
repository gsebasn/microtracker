package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Environment   string
	MongoURI      string
	DatabaseName  string
	ServerAddress string
	RateLimit     RateLimitConfig
}

type RateLimitConfig struct {
	Default   EndpointRateLimit
	Endpoints map[string]EndpointRateLimit
}

type EndpointRateLimit struct {
	RequestsPerMinute int
	BurstSize         int
	TTLMinutes        int
}

// LoadEnv loads the environment file based on APP_ENV
func LoadEnv() error {
	env := getEnv("APP_ENV", "development")
	envFile := fmt.Sprintf(".env.%s", env)

	log.Printf("Loading environment from %s", envFile)

	// Check if the environment file exists
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		log.Printf("Environment file %s not found, using environment variables", envFile)
		return nil
	}

	// Load the environment file
	if err := godotenv.Load(envFile); err != nil {
		return fmt.Errorf("error loading %s: %v", envFile, err)
	}

	return nil
}

func NewConfig() (*Config, error) {
	// Load environment file (optional)
	_ = LoadEnv()

	// Parse default rate limit configuration
	defaultRequestsPerMinute, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS_PER_MINUTE", "100"))
	defaultBurstSize, _ := strconv.Atoi(getEnv("RATE_LIMIT_BURST_SIZE", "50"))
	defaultTTLMinutes, _ := strconv.Atoi(getEnv("RATE_LIMIT_TTL_MINUTES", "5"))

	// Parse endpoint-specific rate limits
	endpoints := make(map[string]EndpointRateLimit)

	// List packages endpoint
	endpoints["GET:/api/v1/packages"] = EndpointRateLimit{
		RequestsPerMinute: getIntEnv("RATE_LIMIT_LIST_REQUESTS_PER_MINUTE", 200),
		BurstSize:         getIntEnv("RATE_LIMIT_LIST_BURST_SIZE", 100),
		TTLMinutes:        getIntEnv("RATE_LIMIT_LIST_TTL_MINUTES", 5),
	}

	// Search packages endpoint
	endpoints["GET:/api/v1/packages/search"] = EndpointRateLimit{
		RequestsPerMinute: getIntEnv("RATE_LIMIT_SEARCH_REQUESTS_PER_MINUTE", 150),
		BurstSize:         getIntEnv("RATE_LIMIT_SEARCH_BURST_SIZE", 75),
		TTLMinutes:        getIntEnv("RATE_LIMIT_SEARCH_TTL_MINUTES", 5),
	}

	// Create package endpoint
	endpoints["POST:/api/v1/packages"] = EndpointRateLimit{
		RequestsPerMinute: getIntEnv("RATE_LIMIT_CREATE_REQUESTS_PER_MINUTE", 50),
		BurstSize:         getIntEnv("RATE_LIMIT_CREATE_BURST_SIZE", 25),
		TTLMinutes:        getIntEnv("RATE_LIMIT_CREATE_TTL_MINUTES", 5),
	}

	config := &Config{
		Environment:   getEnv("APP_ENV", "development"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName:  getEnv("DATABASE_NAME", "tracker"),
		ServerAddress: getEnv("SERVER_ADDRESS", ":9090"),
		RateLimit: RateLimitConfig{
			Default: EndpointRateLimit{
				RequestsPerMinute: defaultRequestsPerMinute,
				BurstSize:         defaultBurstSize,
				TTLMinutes:        defaultTTLMinutes,
			},
			Endpoints: endpoints,
		},
	}

	log.Printf("Loaded configuration: Environment=%s, MongoURI=%s, DatabaseName=%s, ServerAddress=%s",
		config.Environment, config.MongoURI, config.DatabaseName, config.ServerAddress)
	log.Printf("Rate Limit Configuration: Default={RequestsPerMinute=%d, BurstSize=%d, TTLMinutes=%d}",
		config.RateLimit.Default.RequestsPerMinute, config.RateLimit.Default.BurstSize, config.RateLimit.Default.TTLMinutes)
	for endpoint, limit := range config.RateLimit.Endpoints {
		log.Printf("Rate Limit for %s: {RequestsPerMinute=%d, BurstSize=%d, TTLMinutes=%d}",
			endpoint, limit.RequestsPerMinute, limit.BurstSize, limit.TTLMinutes)
	}

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

func ConnectDB(cfg *Config) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return client.Database(cfg.DatabaseName), nil
}
