package api

import (
	"net/http"
	"project/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
	"project/internal/model"
)

type ChatHandler struct {
	chatService service.ChatService
}

func NewChatHandler(svc service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: svc}
}

type CreateChatInput struct {
	InitialMessage string `json:"initialMessage" binding:"required"`
}

type SendMessageInput struct {
	Content string `json:"message" binding:"required"`
}


func (h *ChatHandler) CreateChat(c *gin.Context) {
	var input CreateChatInput
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

	chat, err := h.chatService.CreateChat(UserIDInt, input.InitialMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"chat": chat})
}

func (h *ChatHandler) GetChats(c *gin.Context) {
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

	chats, err := h.chatService.GetChatbyUserID(UserIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chats": chats})
}

func (h *ChatHandler) GetChat(c *gin.Context) {
	chatIDstr := c.Param("id")
	ChatID, err := strconv.ParseUint(chatIDstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
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
	chat, messages, err := h.chatService.GetChatByID(UserIDInt, int(ChatID))
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"chat": chat, "messages": messages})
}


func (h *ChatHandler) SendMessage(c *gin.Context) {
	chatID := c.Param("id")
	ChatIDInt, err := strconv.ParseUint(chatID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var UserIDInt int
	switch v:= userID.(type) {
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

	chat, messages, err := h.chatService.GetChatByID(UserIDInt, int(ChatIDInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var input SendMessageInput
	if err := c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	userMessage := model.Message{
		ChatID:  chat.ID,
		Role:    "user",
		Content: input.Content,
	}

	if err := h.chatService.SaveMessage(int(ChatIDInt), &userMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	messages = append(messages, userMessage)

	aiResponse, err := h.chatService.GenerateAIResponse(messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	aiMessage := model.Message{
		ChatID:  chat.ID,
		Role:    "AI",
		Content: aiResponse,
	}
	if err := h.chatService.SaveMessage(int(ChatIDInt), &aiMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	
	c.JSON(http.StatusOK, gin.H{"chat": chat, "messages": append(messages, aiMessage)})
}