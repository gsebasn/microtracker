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

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting messenger service...")

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
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

// handleSNSWebhook handles incoming SNS notifications
func handleSNSWebhook(c *gin.Context) {
	var message Message
	if err := c.ShouldBindJSON(&message); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Log the received message
	log.Printf("Received message: %+v", message)

	// Handle different message types
	switch message.Type {
	case "SubscriptionConfirmation":
		handleSubscriptionConfirmation(c, message)
	case "Notification":
		handleNotification(c, message)
	default:
		log.Printf("Unknown message type: %s", message.Type)
		c.JSON(http.StatusOK, gin.H{"status": "unknown_message_type"})
	}
}

// handleSubscriptionConfirmation handles SNS subscription confirmation requests
func handleSubscriptionConfirmation(c *gin.Context, message Message) {
	log.Printf("Handling subscription confirmation for topic: %s", message.TopicArn)
	// TODO: Implement subscription confirmation logic
	c.JSON(http.StatusOK, gin.H{"status": "subscription_confirming"})
}

// handleNotification processes regular SNS notifications
func handleNotification(c *gin.Context, message Message) {
	log.Printf("Processing notification from topic: %s", message.TopicArn)

	// Parse the message content
	var content map[string]interface{}
	if err := json.Unmarshal([]byte(message.Message), &content); err != nil {
		log.Printf("Error parsing message content: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message content"})
		return
	}

	// TODO: Implement your message processing logic here
	// For example, you might want to:
	// 1. Validate the message
	// 2. Process the message content
	// 3. Store the message
	// 4. Send acknowledgments

	c.JSON(http.StatusOK, gin.H{"status": "message_received"})
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
