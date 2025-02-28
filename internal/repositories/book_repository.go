package repositories

import (
	"errors"
	"github.com/shani34/book-management-system/internal/models"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetAll(limit, offset int) ([]models.Book, error) {
	var books []models.Book
	result := r.db.Limit(limit).Offset(offset).Find(&books)
	return books, result.Error
}

func (r *BookRepository) GetByID(id uint) (*models.Book, error) {
	var book models.Book
	result := r.db.First(&book, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &book, result.Error
}

func (r *BookRepository) Create(book *models.Book) error {
	result := r.db.Create(book)
	return result.Error
}

func (r *BookRepository) Update(book *models.Book) error {
	result := r.db.Save(book)
	return result.Error
}

func (r *BookRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Book{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}