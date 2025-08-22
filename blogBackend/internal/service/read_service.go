package service

import (
	"project/internal/model"
	"project/internal/store"
	"errors"

)

type ReadService interface {
	CreateReadTime(userID int, read *model.Read) error
	GetWeeklyReadTime(userID int) ([]model.Read, error)
}

type readService struct {
	readTimeStore store.ReadTimeStore
}

func NewReadService(readTimeStore store.ReadTimeStore) ReadService {
	return &readService{readTimeStore: readTimeStore}
}

func (s *readService) CreateReadTime(userID int, read *model.Read) error {
	if read == nil {
		return errors.New("read cannot be nil")
	}
	read.UserID = uint(userID) // Set the UserID for the read time
	return s.readTimeStore.CreateReadTime(userID, read)
}

func (s *readService) GetWeeklyReadTime(userID int) ([]model.Read, error) {
	reads, err := s.readTimeStore.GetWeeklyReadTime(userID)
	if err != nil {
		return nil, err
	}
	return reads, nil
}