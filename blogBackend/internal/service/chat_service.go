package service

import (
	"project/internal/store"
	"project/internal/model"
	"github.com/sashabaranov/go-openai"
	"context"
	"errors"
)

type ChatService interface {
	CreateChat(userID int, initialMessage string) (chat *model.ChatLog, err error)
	SaveMessage(chatID int, message *model.Message) error
	GetChatbyUserID(userID int) (chats []model.ChatLog, err error)
	GetChatByID(chatID int, userID int) (model *model.ChatLog, messages []model.Message, err error)
	GenerateAIResponse(messages []model.Message) (string, error)
	UpdateChat(chat *model.ChatLog) error
}

type chatService struct {
	chatStore    store.ChatLogStore
	messageStore store.MessageStore
	llmClient    *openai.Client
}

func NewChatService(chatStore store.ChatLogStore, messageStore store.MessageStore, API string) ChatService {
	config := openai.DefaultConfig(API)
	config.BaseURL = "https://api.deepseek.ai/v1"
	llmClient := openai.NewClientWithConfig(config)
	return &chatService{
		chatStore:    chatStore,
		messageStore: messageStore,
		llmClient:    llmClient,
	}
}

func (s *chatService) CreateChat(userID int, initialMessage string) (chat *model.ChatLog, err error) {
	chat = &model.ChatLog{
		UserID: uint(userID),
		Title:  "New Chat",
	}
	if err := s.chatStore.CreateChat(chat); err != nil {
		return nil, err
	}

	userMassage := model.Message{
		ChatID:   chat.ID,
		Role:     "AI",
		Content:  initialMessage,
	}

	if err := s.messageStore.CreateMessage(&userMassage); err != nil {
		return nil, err
	}

	aiResponse, err := s.GenerateAIResponse([]model.Message{userMassage})
	if err != nil {
		return nil, err
	}
	

	aiMessage := model.Message{
		ChatID:   chat.ID,
		Role:     "AI",
		Content:  aiResponse,
	}

	if err := s.messageStore.CreateMessage(&aiMessage); err != nil {
		return nil, err
	}

	return chat, nil
}


func (s *chatService) SaveMessage(chatID int, message *model.Message) error {
	return s.messageStore.CreateMessage(message)
}

func (s *chatService) GetChatbyUserID(userID int) (chats []model.ChatLog, err error) {
	return s.chatStore.GetChatByUserID(userID)
}

func (s *chatService) GetChatByID(chatID int, userID int) (model *model.ChatLog, messages []model.Message, err error) {
	chat, err := s.chatStore.GetChatByID(chatID)
	if err != nil {
		return nil, nil, err
	}
	if int(chat.UserID) != userID {
		return nil, nil, errors.New("unauthorized access to chat")
	}
	messages, err = s.messageStore.GetMessagesByChatID(chatID)
	if err != nil {
		return nil, nil, err
	}
	return chat, messages, nil
}


func (s *chatService) GenerateAIResponse(messages []model.Message) (string, error) {
	var openaiMessages []openai.ChatCompletionMessage

	openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
		Role:   "system",
		Content: "You are a helpful assistant.",
	})

	for _, msg := range messages {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	resp, err := s.llmClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    "deepseek-3.5",
			Messages: openaiMessages,
			Temperature: 0.7,
			MaxTokens: 150,
		},
	)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}
	return "", errors.New("no response from AI")
}

func (s *chatService) UpdateChat(chat *model.ChatLog) error {
	return s.chatStore.UpdateChat(chat)
}