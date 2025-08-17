package service

import (
	"errors"
	"project/internal/model"
	"project/internal/store"
)

type UpdateBookLogInput struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	CoverUrl    string `json:"coverUrl"`
	Description string `json:"description"`
	PublishedAt string `json:"publishedAt"`
	ISBN        string `json:"isbn"`
	Category    string `json:"category"`
	Rating      int    `json:"rating"`
	Review      string `json:"review"`
	Status      string `json:"status"`
	MyRating    *int   `json:"myRating"`
	MyComment   string `json:"myComment"`
}

type LogService interface {
	CreateBookLog(userID int, book *model.BookLog) error
	FindBookLogByStatus(userID int, status string) ([]model.BookLog, error)
	GetBookByIDAndUserID(bookID int, userID int) (*model.BookLog, error)
	UpdateLog(BookID int, userID int, params UpdateBookLogInput) (existingLog *model.BookLog, err error)
	SearchBookByTitleOrAuthor(query string) ([]model.BookLog, error)
}

type logService struct {
	bookLogStore store.BookLogStore
}

func NewLogService(bookLogStore store.BookLogStore) LogService {
	return &logService{bookLogStore: bookLogStore}
}

func (s *logService) CreateBookLog(userID int, book *model.BookLog) error {
	if book == nil {
		return errors.New("book cannot be nil")
	}
	book.UserID = uint(userID) // Set the UserID for the book log
	return s.bookLogStore.Create(userID, book)
}

func (s *logService) FindBookLogByStatus(userID int, status string) ([]model.BookLog, error) {
	bookLogs, err := s.bookLogStore.FindBookLogByStatus(userID, status)
	if err != nil {
		return nil, err
	}
	return bookLogs, nil
}

func (s *logService) GetBookByIDAndUserID(bookID int, userID int) (*model.BookLog, error) {
	book, err := s.bookLogStore.GetBookByIDAndUserID(bookID, userID)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *logService) UpdateLog(BookID int, userID int, params UpdateBookLogInput) (existingLog *model.BookLog, err error) {
	existingLog, errorr := s.bookLogStore.GetBookByIDAndUserID(BookID, userID)
	if errorr != nil {
		return nil, errorr
	}
	existingLog.Title = params.Title
	existingLog.Author = params.Author
	existingLog.CoverUrl = params.CoverUrl
	existingLog.Description = params.Description
	existingLog.ISBN = params.ISBN
	existingLog.MyRating = params.MyRating
	existingLog.MyComment = params.MyComment
	existingLog.Status = params.Status

	if err := s.bookLogStore.UpdateLog(existingLog); err != nil {
		return nil, err
	}
	return existingLog, nil
}

func (s *logService) SearchBookByTitleOrAuthor(query string) ([]model.BookLog, error) {
	books, err := s.bookLogStore.SearchBookByTitleOrAuthor(query)
	if err != nil {
		return nil, err
	}
	return books, nil
}