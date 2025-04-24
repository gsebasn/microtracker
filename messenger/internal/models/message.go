package models

import (
	"encoding/json"
	"time"
)

// Message represents a processed SNS message
type Message struct {
	ID          string                 `json:"id"`
	TopicARN    string                 `json:"topic_arn"`
	Content     map[string]interface{} `json:"content"`
	ReceivedAt  time.Time              `json:"received_at"`
	ProcessedAt *time.Time             `json:"processed_at,omitempty"`
	Status      string                 `json:"status"`
	RetryCount  int                    `json:"retry_count"`
	Error       string                 `json:"error,omitempty"`
}

// NewMessage creates a new Message instance
func NewMessage(topicARN string, content map[string]interface{}) *Message {
	return &Message{
		ID:         generateMessageID(),
		TopicARN:   topicARN,
		Content:    content,
		ReceivedAt: time.Now(),
		Status:     "pending",
		RetryCount: 0,
	}
}

// MarkAsProcessed marks the message as processed
func (m *Message) MarkAsProcessed() {
	now := time.Now()
	m.ProcessedAt = &now
	m.Status = "processed"
}

// MarkAsFailed marks the message as failed
func (m *Message) MarkAsFailed(err error) {
	m.Status = "failed"
	m.Error = err.Error()
	m.RetryCount++
}

// IsProcessed returns true if the message has been processed
func (m *Message) IsProcessed() bool {
	return m.Status == "processed"
}

// CanRetry returns true if the message can be retried
func (m *Message) CanRetry(maxRetries int) bool {
	return m.Status == "failed" && m.RetryCount < maxRetries
}

// ToJSON converts the message to JSON
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON creates a Message from JSON
func FromJSON(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// generateMessageID generates a unique message ID
func generateMessageID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
