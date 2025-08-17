package store

import (
	"project/internal/model"

	"gorm.io/gorm"
)

type BookLogStore interface {
	Create(userid int, book *model.BookLog) error
	Migrate() error
	FindBookLogByStatus(userID int, status string) ([]model.BookLog, error)

	//this function used to get the details of a book log by its ID and user ID.
	GetBookByIDAndUserID(bookID int, userID int) (*model.BookLog, error)
	UpdateLog(log *model.BookLog) error
	SearchBookByTitleOrAuthor(query string) ([]model.BookLog, error)
}

type bookLogStore struct {
	db *gorm.DB
}

func NewBookLogStore(db *gorm.DB) BookLogStore {
	return &bookLogStore{db: db}
}

func (s *bookLogStore) Migrate() error {
	return s.db.AutoMigrate(&model.BookLog{})
}

func (s *bookLogStore) Create(userid int, book *model.BookLog) error {
	newbook := &model.BookLog{
		UserID:      uint(userid), // 确保类型匹配
		Title:       book.Title,
		Author:      book.Author,
		CoverUrl:    book.CoverUrl,
		Description: book.Description,
		PublishedAt: book.PublishedAt,
		ISBN:        book.ISBN,
		Category:    book.Category,
		Rating:      book.Rating,
		Review:      book.Review,
		MyRating:    book.MyRating,
		MyComment:   book.MyComment,
		Status:      book.Status,
	}
	return s.db.Create(newbook).Error
}

func (s *bookLogStore) FindBookLogByStatus(userID int, status string) ([]model.BookLog, error) {
	var books []model.BookLog

	query := s.db.Where("user_id = ?", userID)

	if status != "" && status != "All" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil // Assuming you want the first book log found
}

func (s *bookLogStore) GetBookByIDAndUserID(bookID int, userID int) (*model.BookLog, error) {
	var book model.BookLog
	if err := s.db.Where("id = ? AND user_id = ?", bookID, userID).First(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (s *bookLogStore) UpdateLog(log *model.BookLog) error {
	result := s.db.Model(&model.BookLog{}).Where("id = ? AND user_id = ?", log.ID, log.UserID).Updates(log)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *bookLogStore) SearchBookByTitleOrAuthor(query string) ([]model.BookLog, error) {
	var books []model.BookLog
	if err := s.db.Preload("User").Where("title LIKE ? OR author LIKE ?", "%"+query+"%", "%"+query+"%").Limit(20).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}
