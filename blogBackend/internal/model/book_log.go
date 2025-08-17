package model

import (
	"gorm.io/gorm"
)

type BookLog struct {
	gorm.Model
	// Getting from external API
	Title       string `json:"title" gorm:"type:varchar(100)"`
	Author      string `json:"author" gorm:"type:varchar(100)"`
	Description string `json:"description" gorm:"type:text"`
	PublishedAt string `json:"published_at" gorm:"type:varchar(20)"`
	ISBN        string `json:"isbn" gorm:"type:varchar(20);unique"`
	Category    string `json:"category" gorm:"type:varchar(50)"`
	Rating      int    `json:"rating" gorm:"type:int;default:0"`
	Review      string `json:"review" gorm:"type:text"`
	CoverUrl    string `json:"cover_url" gorm:"type:varchar(255)"`

	UserID uint    `json:"user_id" gorm:"not null;index"`
	User   UserLog `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	MyRating  *int   `json:"my_rating" gorm:"type:int;default:0"`
	MyComment string `json:"my_comment" gorm:"type:text"`

	// Book status
	Status string `json:"status" gorm:"type:varchar(20);index"`
}
