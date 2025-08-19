package store

import (
	"project/internal/model"

	"gorm.io/gorm"
)

type ForumStore interface {
	Migrate() error
	CreateTopic(topic *model.Topic) error
	GetTopic(page, pageSize int) ([]model.Topic, int64, error) //this function used to get the topics list with pagination
	GetTopicByID(topicID uint) (*model.Topic, error)           //this function used to get the topic details by its ID
	CreateComment(comment *model.Comment) error
	GetCommentsByTopicID(topicID int, page, pageSize int) ([]model.Comment, error)

	IncrementViewCount(topicID int) error
}

type forumStore struct {
	db *gorm.DB
}

func NewForumStore(db *gorm.DB) ForumStore {
	return &forumStore{db: db}
}

func (s *forumStore) Migrate() error {
	return s.db.AutoMigrate(&model.Topic{}, &model.Comment{})
}

func (s *forumStore) CreateTopic(topic *model.Topic) error {
	return s.db.Create(topic).Error
}

func (s *forumStore) GetTopic(page, pageSize int) ([]model.Topic, int64, error) {
	var topics []model.Topic
	var total int64

	offset := (page - 1) * pageSize

	//count total topics
	if err := s.db.Model(&model.Topic{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	//page and pageSize
	//page 1 pageSize 20 means the first page will show 20 topics
	//and offset means the number of topics to skip

	if err := s.db.Preload("User").Offset(offset).Limit(pageSize).Find(&topics).Error; err != nil {
		return nil, 0, err
	}
	return topics, total, nil
}

func (s *forumStore) GetTopicByID(topicID uint) (*model.Topic, error) {
	var topic model.Topic
	if err := s.db.Preload("User").Preload("Comments").Preload("Comments.User").First(&topic).Error; err != nil {
		return nil, err
	}
	return &topic, nil
}

func (s *forumStore) CreateComment(comment *model.Comment) error {
	return s.db.Create(comment).Error
}

func (s *forumStore) GetCommentsByTopicID(topicID int, page, pageSize int) ([]model.Comment, error) {
	var comments []model.Comment
	var total int64
	// Count total comments for the topic
	if err := s.db.Model(&model.Comment{}).Where("topic_id = ?", topicID).Count(&total).Error; err != nil {
		return nil, err
	}
	offset := (page - 1) * pageSize
	if err := s.db.Preload("User").Order("created_at ASC").Where("topic_id = ?", topicID).Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *forumStore) IncrementViewCount(topicID int) error {
	return s.db.Model(&model.Topic{}).Where("id = ?", topicID).Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}
