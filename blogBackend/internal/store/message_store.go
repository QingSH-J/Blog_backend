package store

import (
	"project/internal/model"
	"gorm.io/gorm"
)

type MessageStore interface {
	Migrate() error
	CreateMessage(message *model.Message) error
	GetMessagesByChatID(chatID int) ([]model.Message, error)
}

type messageStore struct {
	db *gorm.DB
}

func NewMessageStore(db *gorm.DB) MessageStore {
	return &messageStore{db: db}
}

func (s *messageStore) Migrate() error {
	return s.db.AutoMigrate(&model.Message{})
}

func (s *messageStore) CreateMessage(message *model.Message) error {
	return s.db.Create(message).Error
}

func (s *messageStore) GetMessagesByChatID(chatID int) ([]model.Message, error){
	var messages []model.Message
	if err := s.db.Where("chat_id = ?", chatID).Order("created_at ASC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}