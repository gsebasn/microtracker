package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Message represents the structure of an SNS message
type Message struct {
	Type             string `json:"Type"`
	MessageId        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
}

var logger *zap.Logger

func init() {
	// Configure Zap logger
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var err error
	logger, err = config.Build()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()
}

func main() {
	// Set up logging
	logger.Info("Starting messenger service...")

	// Create Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// SNS webhook endpoint
	router.POST("/webhook", handleSNSWebhook)

	// Start server
	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", zap.String("port", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
}

// handleSNSWebhook handles incoming SNS notifications
func handleSNSWebhook(c *gin.Context) {
	var message Message
	if err := c.ShouldBindJSON(&message); err != nil {
		logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Log the received message
	logger.Info("Received message",
		zap.String("message_id", message.MessageId),
		zap.String("topic_arn", message.TopicArn),
		zap.String("type", message.Type),
	)

	// Handle different message types
	switch message.Type {
	case "SubscriptionConfirmation":
		handleSubscriptionConfirmation(c, message)
	case "Notification":
		handleNotification(c, message)
	default:
		logger.Warn("Unknown message type", zap.String("type", message.Type))
		c.JSON(http.StatusOK, gin.H{"status": "unknown_message_type"})
	}
}

// handleSubscriptionConfirmation handles SNS subscription confirmation requests
func handleSubscriptionConfirmation(c *gin.Context, message Message) {
	logger.Info("Handling subscription confirmation",
		zap.String("topic_arn", message.TopicArn),
	)
	// TODO: Implement subscription confirmation logic
	c.JSON(http.StatusOK, gin.H{"status": "subscription_confirming"})
}

// handleNotification processes regular SNS notifications
func handleNotification(c *gin.Context, message Message) {
	logger.Info("Processing notification",
		zap.String("topic_arn", message.TopicArn),
	)

	// Parse the message content
	var content map[string]interface{}
	if err := json.Unmarshal([]byte(message.Message), &content); err != nil {
		logger.Error("Error parsing message content",
			zap.Error(err),
			zap.String("message_id", message.MessageId),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message content"})
		return
	}

	// TODO: Implement your message processing logic here
	// For example, you might want to:
	// 1. Validate the message
	// 2. Process the message content
	// 3. Store the message
	// 4. Send acknowledgments

	logger.Info("Message processed successfully",
		zap.String("message_id", message.MessageId),
		zap.String("topic_arn", message.TopicArn),
	)
	c.JSON(http.StatusOK, gin.H{"status": "message_received"})
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
