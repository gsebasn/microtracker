package config

import (
	"context"
	"fmt"
	"log"
	"os"
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
}

// LoadEnv loads the environment file based on APP_ENV
func LoadEnv() error {
	env := getEnv("APP_ENV", "development")
	envFile := fmt.Sprintf(".env.%s", env)

	log.Printf("Loading environment from %s", envFile)

	// Check if the environment file exists
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return fmt.Errorf("environment file %s not found", envFile)
	}

	// Load the environment file
	if err := godotenv.Load(envFile); err != nil {
		return fmt.Errorf("error loading %s: %v", envFile, err)
	}

	return nil
}

func NewConfig() (*Config, error) {
	// Load environment file
	if err := LoadEnv(); err != nil {
		return nil, fmt.Errorf("failed to load environment: %v", err)
	}

	config := &Config{
		Environment:   getEnv("APP_ENV", "development"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName:  getEnv("DATABASE_NAME", "tracker"),
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
	}

	log.Printf("Loaded configuration: Environment=%s, MongoURI=%s, DatabaseName=%s, ServerAddress=%s",
		config.Environment, config.MongoURI, config.DatabaseName, config.ServerAddress)

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
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
