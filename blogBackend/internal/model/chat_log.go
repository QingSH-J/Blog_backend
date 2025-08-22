package model

import (
	"gorm.io/gorm"
)

type ChatLog struct {
	gorm.Model
	UserID   uint      `json:"user_id" gorm:"not null;index"`
	User     UserLog   `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Title    string    `json:"title" gorm:"type:varchar(100)"`
	LastActive string  `json:"last_activity" gorm:"index"`
	Messages []Message `json:"messages" gorm:"foreignKey:ChatLogID;"`
}

type Message struct {
	gorm.Model
	ChatID    uint   `json:"chat_id" gorm:"not null;index"`
	Chat      ChatLog `json:"chat" gorm:"foreignKey:ChatID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ChatLogID uint   `json:"chat_log_id" gorm:"not null;index"`
	Role      string `json:"role" gorm:"type:varchar(20)"`
	Content   string `json:"content" gorm:"type:text"`
}
