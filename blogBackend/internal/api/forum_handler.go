package api

import (
	"net/http"
	"project/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ForumHandler struct {
	forumService service.ForumService
}

func NewForumHandler(svc service.ForumService) *ForumHandler {
	return &ForumHandler{forumService: svc}
}

type CreateTopicInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type CreateCommentInput struct {
	Content string `json:"content" binding:"required"`
}

func (h *ForumHandler) CreateTopic(c *gin.Context) {
	var input CreateTopicInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found"})
		return
	}

	var userIDUint uint

	switch v := userID.(type) {
	case uint:
		userIDUint = v
	case int:
		userIDUint = uint(v)
	case float64:
		userIDUint = uint(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}

	topic, err := h.forumService.CreateTopic(userIDUint, input.Title, input.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"topic": topic})
}

func (h *ForumHandler) GetTopics(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	topics, total, err := h.forumService.GetTopic(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"topics": topics, "page": gin.H{"current": page, "size": pageSize, "total": total, "totalPages": (total + int64(pageSize) - 1) / int64(pageSize)}})
}

func (h *ForumHandler) GetTopicByID(c *gin.Context) {
	topicIDStr := c.Param("id")
	topicID, err := strconv.ParseUint(topicIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}
	topic, err := h.forumService.GetTopicByID(uint(topicID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if topic == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
		return
	}

	go h.forumService.IncrementViewCount(int(topicID))
	c.JSON(http.StatusOK, gin.H{"topic": topic})
}

func (h *ForumHandler) CreateComment(c *gin.Context) {
	var input CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user ID not found"})
		return
	}

	var UserIDUint uint
	switch v := userID.(type) {
	case uint:
		UserIDUint = v
	case int:
		UserIDUint = uint(v)
	case float64:
		UserIDUint = uint(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}

	topicIDStr := c.Param("id")
	topicID, err := strconv.ParseUint(topicIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	comment, err := h.forumService.CreateComment(UserIDUint, int(topicID), input.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"comment": comment})
}

func (h *ForumHandler) GetComments(c *gin.Context) {
	topicIDStr := c.Param("id")
	topicID, err := strconv.ParseUint(topicIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	comments, err := h.forumService.GetCommentsByTopicID(int(topicID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"comments": comments, "page": gin.H{"current": page, "size": pageSize, "total": len(comments), "totalPages": (len(comments) + pageSize - 1) / pageSize}})
}
