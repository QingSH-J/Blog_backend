package api

import (
	"log"
	"net/http"
	"project/internal/model"
	"project/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	logService service.LogService
}

func NewLogHandler(svc service.LogService) *LogHandler {
	return &LogHandler{logService: svc}
}

type UpdateLogInput struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	CoverUrl    string `json:"coverUrl"`
	Description string `json:"description"`
	PublishedAt string `json:"publishedAt"`
	ISBN        string `json:"isbn"`
	Category    string `json:"category"`
	Rating      int    `json:"rating"`
	Review      string `json:"review"`
	Status      string `json:"status"`
	MyRating    *int   `json:"myRating"`
	MyComment   string `json:"myComment"`
}

type CreateBookLogInput struct {
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	CoverUrl    string `json:"coverUrl"`
	Description string `json:"description"`
	PublishedAt string `json:"publishedAt"`
	ISBN        string `json:"isbn"`
	Category    string `json:"category"`
	Rating      int    `json:"rating"`
	Review      string `json:"review"`
	Status      string `json:"status" binding:"required"`
	MyRating    *int   `json:"myRating"`
	MyComment   string `json:"myComment"`
}

func (h *LogHandler) CreateBookLog(c *gin.Context) {
	var input CreateBookLogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Handle different possible types for userID
	var userIDInt int
	switch v := userID.(type) {
	case float64:
		userIDInt = int(v)
	case int:
		userIDInt = v
	case uint:
		userIDInt = int(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID format"})
		return
	}

	bookLog := &model.BookLog{
		Title:       input.Title,
		Author:      input.Author,
		CoverUrl:    input.CoverUrl,
		Description: input.Description,
		PublishedAt: input.PublishedAt,
		ISBN:        input.ISBN,
		Category:    input.Category,
		Rating:      input.Rating,
		Review:      input.Review,
		Status:      input.Status,
		MyRating:    input.MyRating,
		MyComment:   input.MyComment,
	}
	if err := h.logService.CreateBookLog(userIDInt, bookLog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Book log created successfully"})
}

func (h *LogHandler) GetBookLog(c *gin.Context) {
	// 安全地获取并转换userID
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// 安全类型转换
	var userIDInt int
	switch v := userIDRaw.(type) {
	case float64:
		userIDInt = int(v)
	case int:
		userIDInt = v
	case uint:
		userIDInt = int(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID format"})
		return
	}

	status := c.Query("status")

	books, err := h.logService.FindBookLogByStatus(userIDInt, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(books) == 0 {
		c.JSON(http.StatusOK, gin.H{"books": []interface{}{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books})
}

func (h *LogHandler) GetBook(c *gin.Context) {
	idStr := c.Param("id")
	BookID, err := strconv.ParseUint(idStr, 10, 64)
	log.Println("BookID:", BookID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	var userIDInt int
	switch v := userIDRaw.(type) {
	case float64:
		userIDInt = int(v)
	case int:
		userIDInt = v
	case uint:
		userIDInt = int(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID format"})
		return
	}
	book, err := h.logService.GetBookByIDAndUserID(int(BookID), userIDInt)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": book})
}

func (h *LogHandler) UpdateBookLog(c *gin.Context) {
	idStr := c.Param("id")
	BookID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	var userIDInt int
	switch v := userIDRaw.(type) {
	case float64:
		userIDInt = int(v)
	case int:
		userIDInt = v
	case uint:
		userIDInt = int(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID format"})
		return
	}

	var params service.UpdateBookLogInput
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatelog, err := h.logService.UpdateLog(int(BookID), userIDInt, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book log updated successfully", "book": updatelog})
}

func (h *LogHandler) SearchBook(c *gin.Context) {
	q := c.Query("query")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}
	books, err := h.logService.SearchBookByTitleOrAuthor(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(books) == 0 {
		c.JSON(http.StatusOK, gin.H{"books": []interface{}{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books})
}
