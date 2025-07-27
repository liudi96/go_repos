package handler

import (
	"fmt"
	"net/http"
	"palu-wiki-backend/internal/models"
	"palu-wiki-backend/internal/repository"
	"palu-wiki-backend/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	guideRepo     *repository.GuideRepository
	updateService *service.UpdateService
}

func NewAdminHandler(guideRepo *repository.GuideRepository, updateService *service.UpdateService) *AdminHandler {
	return &AdminHandler{
		guideRepo:     guideRepo,
		updateService: updateService,
	}
}

type CreateGuideTopicRequest struct {
	Topic string `json:"topic" binding:"required"`
}

// CreateGuideTopic handles the creation of a new guide topic from the admin backend.
// This topic will then be used by the AI to generate a full guide.
func (h *AdminHandler) CreateGuideTopic(c *gin.Context) {
	var req CreateGuideTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Here, we'll trigger the AI to generate a guide based on the topic.
	// For simplicity, we'll call a method on updateService.
	// In a real-world scenario, this might be an asynchronous task.
	dummyUpdate := models.OfficialUpdate{
		Title:       fmt.Sprintf("%s", req.Topic),
		Content:     fmt.Sprintf("根据主题“%s”生成的攻略。", req.Topic),
		SourceURL:   fmt.Sprintf("admin-topic-%s-%d", req.Topic, time.Now().UnixNano()), // Generate a unique URL
		PublishDate: time.Now(),
	}

	// Call the update service to process this "update" (which is our new topic)
	// This will trigger AI generation and storage.
	err := h.updateService.ProcessUpdates([]models.OfficialUpdate{dummyUpdate})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate guide for topic: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Guide generation for topic initiated successfully."})
}

// GetGuidesForAdmin retrieves all guides, potentially with more details for admin view.
func (h *AdminHandler) GetGuidesForAdmin(c *gin.Context) {
	guides, err := h.guideRepo.GetAllGuides()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch guides: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": guides})
}
