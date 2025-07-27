package models

import (
	"encoding/json"
	"time"
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

// UserQuery represents a user's query to the AI bot.
type UserQuery struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     string    `gorm:"type:varchar(100)" json:"user_id"` // WeChat OpenID or UnionID
	QueryText  string    `gorm:"type:text;not null" json:"query_text"`
	AIResponse string    `gorm:"type:longtext;not null" json:"ai_response"`
	QueryTime  time.Time `json:"query_time"`
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
