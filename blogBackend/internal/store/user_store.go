package store

import (
	"project/internal/model"
	"gorm.io/gorm"

)

type UserStore interface {
	CreateUser(user *model.UserLog) error
	Migrate() error
	FindUserByEmail(email string) (*model.UserLog, error)
}

type userStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) UserStore {
	return &userStore{db: db}
}

func (s *userStore) Migrate() error {
	return s.db.AutoMigrate(&model.UserLog{})
}

func (s *userStore) CreateUser(user *model.UserLog) error {
	return s.db.Create(user).Error
}
func (s *userStore) FindUserByEmail(email string) (*model.UserLog, error) {
	var user model.UserLog
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
