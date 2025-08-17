package service

import (
	"project/internal/model"
	"project/internal/store"
)

type ForumService interface {
	CreateTopic(userID uint, title string, content string) (*model.Topic, error)
	GetTopic(page, pageSize int) ([]model.Topic, int64, error)
	GetTopicByID(topicID uint) (*model.Topic, error)
	CreateComment(userID uint, topicID int, content string) (*model.Comment, error)
	GetCommentsByTopicID(topicID int, page, pageSize int) ([]model.Comment,error)
	IncrementViewCount(topicID int) error
}

type forumService struct {
	forumStore store.ForumStore
}

func NewForumService(forumStore store.ForumStore) ForumService {
	return &forumService{forumStore: forumStore}
}

func (s *forumService) CreateTopic(userID uint, title string, content string) (*model.Topic, error) {
	topic := &model.Topic{
		Title:   title,
		Content: content,
		UserID:  userID,
	}
	if err := s.forumStore.CreateTopic(topic); err != nil {
		return nil, err
	}
	return topic, nil
}

func (s *forumService) GetTopic(page, pageSize int) ([]model.Topic, int64, error) {
	return s.forumStore.GetTopic(page, pageSize)
}

func (s *forumService) GetTopicByID(topicID uint) (*model.Topic, error) {
	topic, err := s.forumStore.GetTopicByID(topicID)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func (s *forumService) CreateComment(userID uint, topicID int, content string) (*model.Comment, error) {
	comment := &model.Comment{
		Content: content,
		TopicID: uint(topicID),
		UserID:  userID,
	}
	if err := s.forumStore.CreateComment(comment); err != nil{
		return nil, err
	}
	return comment, nil
}

func (s *forumService) GetCommentsByTopicID(topicID int, page int, pageSize int) ([]model.Comment, error){
	return s.forumStore.GetCommentsByTopicID(topicID, page, pageSize)
}

func (s *forumService) IncrementViewCount(topicID int) error {
	return s.forumStore.IncrementViewCount(topicID)
}
