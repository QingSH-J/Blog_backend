package model

// Forum model includes Topic and Comment structures

import (
	"gorm.io/gorm"
)

type Topic struct {
	gorm.Model
	Title    string    `json:"title" gorm:"type:varchar(100);not null"`
	Content  string    `json:"content" gorm:"type:text;not null"`
	UserID   uint      `json:"user_id" gorm:"not null;index"`
	User     UserLog   `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comments []Comment `json:"comments" gorm:"foreignKey:TopicID"`
	Viewcount int      `json:"view_count" gorm:"default:0"`
}

type Comment struct {
	gorm.Model
	Content string  `json:"content" gorm:"type:text;not null"`
	TopicID uint    `json:"topic_id" gorm:"not null;index"`
	UserID  uint    `json:"user_id" gorm:"not null;index"`
	User    UserLog `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
