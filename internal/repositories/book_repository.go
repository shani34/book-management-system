package repositories

import (
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
	return &book, result.Error
}

// Implement other CRUD operations...