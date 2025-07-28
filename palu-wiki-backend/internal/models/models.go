package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Guide represents a game guide entry.
type Guide struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Title         string    `gorm:"type:varchar(255);not null" json:"title"`
	Content       string    `gorm:"type:longtext;not null" json:"content"`
	Tags          string    `gorm:"type:varchar(255)" json:"tags"` // Comma-separated tags
	Category      string    `gorm:"type:varchar(100)" json:"category"`
	SourceURL     string    `gorm:"type:varchar(255)" json:"source_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	IsAIGenerated bool      `json:"is_ai_generated"`
	Version       string    `gorm:"type:varchar(50)" json:"version"`
}

// MarshalJSON implements the json.Marshaler interface for Guide.
func (g Guide) MarshalJSON() ([]byte, error) {
	type Alias Guide // Create an alias to avoid infinite recursion

	return json.Marshal(&struct {
		Alias
		CreatedAt string `json:"created_at"` // Override CreatedAt with string type
		UpdatedAt string `json:"updated_at"` // Override UpdatedAt with string type
	}{
		Alias:     (Alias)(g),
		CreatedAt: g.CreatedAt.Format("2006年01月02日"), // Format to YYYY年MM月DD日
		UpdatedAt: g.UpdatedAt.Format("2006年01月02日"), // Format to YYYY年MM月DD日
	})
}

// ChatMessage represents a single message in a chat conversation.
type ChatMessage struct {
	Role    string `json:"role"`    // "user" or "ai"
	Content string `json:"content"` // Message content
}

// UserQuery represents a user's query to the AI bot, including the full conversation history.
type UserQuery struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"type:varchar(100)" json:"user_id"` // WeChat OpenID or UnionID
	Messages  string    `gorm:"type:longtext;not null" json:"-"`  // Store JSON string of ChatMessage array
	QueryTime time.Time `json:"query_time"`
}

// BeforeSave hook to marshal messages to JSON string
func (uq *UserQuery) BeforeSave(tx *gorm.DB) (err error) {
	if len(uq.Messages) == 0 {
		// If Messages is empty, it means it's not set from the frontend,
		// so we don't need to marshal it. This can happen when retrieving from DB.
		return nil
	}
	messagesJSON, err := json.Marshal(uq.Messages)
	if err != nil {
		return err
	}
	uq.Messages = string(messagesJSON)
	return nil
}

// AfterFind hook to unmarshal messages from JSON string
func (uq *UserQuery) AfterFind(tx *gorm.DB) (err error) {
	if uq.Messages == "" {
		return nil
	}
	var chatMessages []ChatMessage
	err = json.Unmarshal([]byte(uq.Messages), &chatMessages)
	if err != nil {
		return err
	}
	// Assign the unmarshaled messages back to a transient field if needed,
	// or handle it directly in the handler. For now, we'll just unmarshal.
	return nil
}

// OfficialUpdate represents an official game update log.
type OfficialUpdate struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Content     string    `gorm:"type:longtext;not null" json:"content"`
	PublishDate time.Time `json:"publish_date"`
	SourceURL   string    `gorm:"type:varchar(255)" json:"source_url"`
	Processed   bool      `gorm:"default:false" json:"processed"`
}
