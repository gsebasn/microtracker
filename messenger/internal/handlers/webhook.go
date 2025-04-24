package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snavarro/microtracker/messenger/internal/models"
)

// WebhookHandler handles SNS webhook requests
type WebhookHandler struct {
	config *Config
}

// Config holds the handler configuration
type Config struct {
	MaxRetries    int
	RetryInterval time.Duration
}

// NewWebhookHandler creates a new WebhookHandler
func NewWebhookHandler(config *Config) *WebhookHandler {
	return &WebhookHandler{
		config: config,
	}
}

// HandleSNSWebhook handles incoming SNS notifications
func (h *WebhookHandler) HandleSNSWebhook(c *gin.Context) {
	var snsMessage struct {
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

	if err := c.ShouldBindJSON(&snsMessage); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Handle different message types
	switch snsMessage.Type {
	case "SubscriptionConfirmation":
		h.handleSubscriptionConfirmation(c, snsMessage)
	case "Notification":
		h.handleNotification(c, snsMessage)
	default:
		log.Printf("Unknown message type: %s", snsMessage.Type)
		c.JSON(http.StatusOK, gin.H{"status": "unknown_message_type"})
	}
}

// handleSubscriptionConfirmation handles SNS subscription confirmation requests
func (h *WebhookHandler) handleSubscriptionConfirmation(c *gin.Context, message struct {
	Type             string `json:"Type"`
	MessageId        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
}) {
	log.Printf("Handling subscription confirmation for topic: %s", message.TopicArn)
	// TODO: Implement subscription confirmation logic
	c.JSON(http.StatusOK, gin.H{"status": "subscription_confirming"})
}

// handleNotification processes regular SNS notifications
func (h *WebhookHandler) handleNotification(c *gin.Context, message struct {
	Type             string `json:"Type"`
	MessageId        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
}) {
	log.Printf("Processing notification from topic: %s", message.TopicArn)

	// Parse the message content
	var content map[string]interface{}
	if err := json.Unmarshal([]byte(message.Message), &content); err != nil {
		log.Printf("Error parsing message content: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message content"})
		return
	}

	// Create a new message
	msg := models.NewMessage(message.TopicArn, content)

	// Process the message
	if err := h.processMessage(msg); err != nil {
		log.Printf("Error processing message: %v", err)
		msg.MarkAsFailed(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process message"})
		return
	}

	msg.MarkAsProcessed()
	c.JSON(http.StatusOK, gin.H{"status": "message_received"})
}

// processMessage processes a message
func (h *WebhookHandler) processMessage(msg *models.Message) error {
	// TODO: Implement your message processing logic here
	// For example:
	// 1. Validate the message content
	// 2. Store the message in a database
	// 3. Send notifications
	// 4. Update related systems

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)
	return nil
}
