package model
import (
	"gorm.io/gorm"
)
type Read struct {
	gorm.Model
	Time int    `json:"time" gorm:"not null"`
	UserID  uint  `json:"user_id" gorm:"not null;index"`
	User    UserLog `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

