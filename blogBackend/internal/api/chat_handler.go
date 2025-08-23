package api

import (
	"log"
	"net/http"
	"project/internal/model"
	"project/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
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
	chat, messages, err := h.chatService.GetChatByID(int(ChatID), UserIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"chat": chat, "messages": messages})
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	// from the URL, get the chat ID
	chatIDStr := c.Param("id")
	chatID, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID format"})
		return
	}

	// record the requested chat ID
	log.Printf("SendMessage - Request for chat ID: %d", chatID)

	// get user ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var userIDInt int
	switch v := userID.(type) {
	case int:
		userIDInt = v
	case uint:
		userIDInt = int(v)
	case float64:
		userIDInt = int(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID type"})
		return
	}

	// get the message content from the request body
	var input SendMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 
	chat, messages, err := h.chatService.GetChatByID(int(chatID), userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// create user message
	userMessage := model.Message{
		ChatID:    chat.ID,
		ChatLogID: chat.ID, // set the ChatLogID same as ChatID to resolve foreign key constraint issue
		Role:      "user",
		Content:   input.Content,
	}

	if err := h.chatService.SaveMessage(int(chatID), &userMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	messages = append(messages, userMessage)

	// generate AI response
	aiResponse, err := h.chatService.GenerateAIResponse(messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// save AI message
	aiMessage := model.Message{
		ChatID:    chat.ID,
		ChatLogID: chat.ID, // set the ChatLogID same as ChatID to resolve foreign key constraint issue
		Role:      "AI",
		Content:   aiResponse,
	}

	if err := h.chatService.SaveMessage(int(chatID), &aiMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat": chat, "messages": append(messages, aiMessage)})
}
