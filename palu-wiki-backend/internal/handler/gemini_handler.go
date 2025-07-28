package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"palu-wiki-backend/internal/models"
	"palu-wiki-backend/internal/repository"
	"palu-wiki-backend/pkg/gemini"

	"github.com/gin-gonic/gin"
)

type GeminiHandler struct {
	geminiClient *gemini.Client
	chatRepo     *repository.ChatRepository
}

func NewGeminiHandler(client *gemini.Client, chatRepo *repository.ChatRepository) *GeminiHandler {
	return &GeminiHandler{geminiClient: client, chatRepo: chatRepo}
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GenerateContentRequest struct {
	Messages []ChatMessage `json:"messages" binding:"required"`
}

type GenerateContentResponse struct {
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

func (h *GeminiHandler) GenerateContent(c *gin.Context) {
	var req GenerateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GenerateContentResponse{Error: err.Error()})
		return
	}

	// Convert handler's ChatMessage to pkg/gemini's ChatMessage
	var geminiMessages []gemini.ChatMessage
	for _, msg := range req.Messages {
		geminiMessages = append(geminiMessages, gemini.ChatMessage{Role: msg.Role, Content: msg.Content})
	}

	content, err := h.geminiClient.GenerateContent(c.Request.Context(), geminiMessages)
	if err != nil {
		log.Printf("Gemini API call failed: %v", err) // Add detailed error logging
		c.JSON(http.StatusInternalServerError, GenerateContentResponse{Error: "AI服务内部错误，请检查后端日志。" + err.Error()})
		return
	}

	// Truncate content if it's too long to avoid potential network issues with large responses
	// WeChat Mini Program might also have implicit limits on response size.
	maxContentLength := 2000 // Example: Limit to 2000 characters
	if len(content) > maxContentLength {
		content = content[:maxContentLength] + "..." // Add ellipsis
	}

	// Save user query and AI response to database
	aiMessage := ChatMessage{Role: "ai", Content: content}
	updatedMessages := append(req.Messages, aiMessage)
	updatedMessagesJSON, err := json.Marshal(updatedMessages)
	if err != nil {
		log.Printf("Failed to marshal updated messages: %v", err)
		c.JSON(http.StatusInternalServerError, GenerateContentResponse{Error: "内部错误：无法处理聊天记录。"})
		return
	}

	userQuery := models.UserQuery{
		// UserID needs to be obtained from authentication context, for now, use a placeholder
		UserID:    "test_user_id",
		Messages:  string(updatedMessagesJSON),
		QueryTime: time.Now(),
	}

	if err := h.chatRepo.SaveUserQuery(&userQuery); err != nil {
		log.Printf("Failed to save user query: %v", err)
		// Do not return error to client, as AI response was successful
	}

	c.JSON(http.StatusOK, GenerateContentResponse{Content: content})
}

type GetChatHistoryResponse struct {
	History []ChatMessage `json:"history"`
	Error   string        `json:"error,omitempty"`
}

func (h *GeminiHandler) GetChatHistory(c *gin.Context) {
	userID := c.Query("user_id") // Get user ID from query parameter, replace with actual auth later
	if userID == "" {
		c.JSON(http.StatusBadRequest, GetChatHistoryResponse{Error: "User ID is required."})
		return
	}

	history, err := h.chatRepo.GetChatHistory(userID)
	if err != nil {
		log.Printf("Failed to get chat history: %v", err)
		c.JSON(http.StatusInternalServerError, GetChatHistoryResponse{Error: "获取聊天历史失败。"})
		return
	}

	var chatMessages []ChatMessage
	for _, entry := range history {
		var messagesInEntry []ChatMessage
		if err := json.Unmarshal([]byte(entry.Messages), &messagesInEntry); err != nil {
			log.Printf("Failed to unmarshal messages from history entry: %v", err)
			continue // Skip this entry if unmarshalling fails
		}
		chatMessages = append(chatMessages, messagesInEntry...)
	}

	c.JSON(http.StatusOK, GetChatHistoryResponse{History: chatMessages})
}
