package repository

import (
	"palu-wiki-backend/internal/models"

	"gorm.io/gorm"
)

// ChatRepository handles database operations for chat messages.
type ChatRepository struct {
	db *gorm.DB
}

// NewChatRepository creates a new ChatRepository.
func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// SaveUserQuery saves a user query and AI response to the database.
func (r *ChatRepository) SaveUserQuery(query *models.UserQuery) error {
	return r.db.Create(query).Error
}

// GetChatHistory retrieves chat history for a given user ID.
// It returns a list of UserQuery objects, ordered by query time.
func (r *ChatRepository) GetChatHistory(userID string) ([]models.UserQuery, error) {
	var history []models.UserQuery
	err := r.db.Where("user_id = ?", userID).Order("query_time ASC").Find(&history).Error
	return history, err
}
