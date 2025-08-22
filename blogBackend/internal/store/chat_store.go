package store

import (
	"project/internal/model"
	"gorm.io/gorm"
)

type ChatLogStore interface {
	Migrate() error
	CreateChat(chat *model.ChatLog) error
	UpdateChat(chat *model.ChatLog) error
	GetChatByUserID(userID int) ([]model.ChatLog, error)
	GetChatByID(chatID int) (*model.ChatLog, error)
}

type chatLogStore struct {
	db *gorm.DB
}

func NewChatLogStore(db *gorm.DB) ChatLogStore {
	return &chatLogStore{db: db}
}

func (s *chatLogStore) Migrate() error {
	return s.db.AutoMigrate(&model.ChatLog{})
}

func (s *chatLogStore) CreateChat(chat *model.ChatLog) error {
	return s.db.Create(chat).Error
}

func (s *chatLogStore) UpdateChat(chat *model.ChatLog) error{
	return s.db.Model(&model.ChatLog{}).Where("id = ?", chat.ID).Updates(chat).Error
}

func (s *chatLogStore) GetChatByUserID(userID int) ([]model.ChatLog, error) {
	var chats []model.ChatLog
	if err := s.db.Preload("User").Where("user_id = ?", userID).Order("updated_at DESC").Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

func (s *chatLogStore) GetChatByID(chatID int) (*model.ChatLog, error) {
	var chat model.ChatLog
	if err := s.db.Preload("User").Where("id = ?", chatID).First(&chat).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}