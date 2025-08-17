package model

import (
	"gorm.io/gorm"
)

type UserLog struct {
	gorm.Model
	UserName string `json:"user_name" gorm:"type:varchar(50);not null"`
	Email    string `json:"email" gorm:"type:varchar(100);unique;not null"`
	Password string `json:"password" gorm:"type:varchar(100);not null"`

	// 反向关联 - 用户添加的所有图书
	BookLogs []BookLog `json:"book_logs" gorm:"foreignKey:UserID"`
}
