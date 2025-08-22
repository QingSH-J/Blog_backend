package store

import (
	"project/internal/model"
	"time"
	"gorm.io/gorm"
)

type ReadTimeStore interface {
	CreateReadTime(userID int, read *model.Read) error
	Migrate() error
	// GetReadTimeByUserID(userID int) (*model.Read, error)
	// UpdateReadTime(userID int, read *model.Read) error
	GetWeeklyReadTime(userID int) ([]model.Read, error)
}

type readTimeStore struct {
	db *gorm.DB
}

func NewReadTimeStore(db *gorm.DB) ReadTimeStore {
	return &readTimeStore{db: db}
}

func (s *readTimeStore) Migrate() error {
	return s.db.AutoMigrate(&model.Read{})
}

func (s *readTimeStore) CreateReadTime(userID int, read *model.Read) error {
	readtime := &model.Read{
		UserID: uint(userID), // 确保类型匹配
		Time:   read.Time,
	}
	return s.db.Create(readtime).Error
}

func (s *readTimeStore) GetWeeklyReadTime(userID int) ([]model.Read, error) {
	var reads []model.Read
	timeRange := time.Now().AddDate(0, 0, -7) // get the last 7 days
	if err := s.db.Where("user_id = ? AND created_at >= ?", userID, timeRange).Order("created_at ASC").Find(&reads).Error; err != nil {
		return nil, err
	}
	return reads, nil
}