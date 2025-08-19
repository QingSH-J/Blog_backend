package api

import (
	"net/http"
	"project/internal/model"
	"project/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

type ReadInput struct {
	Time int       `json:"data" binding:"required"`
	Date time.Time `json:"date" binding:"required"`
}
type ReadHandler struct {
	CreateReadService service.ReadService
}

func NewReadHandler(svc service.ReadService) *ReadHandler {
	return &ReadHandler{CreateReadService: svc}
}

func (h *ReadHandler) CreateReadTime(c *gin.Context) {
	var input ReadInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var UserIDInt int
	switch v := userID.(type) {
	case int:
		UserIDInt = v
	case uint:
		UserIDInt = int(v)
	case float64:
		UserIDInt = int(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}
	read := &model.Read{
		Time: input.Time,
	}

	err := h.CreateReadService.CreateReadTime(UserIDInt, read)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Read time created successfully"})
}
