package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shani34/book-management-system/internal/models"
	"github.com/shani34/book-management-system/internal/repositories"
	"github.com/shani34/book-management-system/pkg/kafka"
	"github.com/shani34/book-management-system/pkg/redis"
	"log"
	"time"
)

type BookService struct {
	repo    *repositories.BookRepository
	cache   *redis.Client
	timeout time.Duration
}

func NewBookService(repo *repositories.BookRepository, cache *redis.Client) *BookService {
	return &BookService{
		repo:    repo,
		cache:   cache,
		timeout: 10 * time.Minute,
	}
}

func (s *BookService) GetAllBooks(limit, offset int) ([]models.Book, error) {
	cacheKey := fmt.Sprintf("books:%d:%d", limit, offset)
	
	// Try cache first
	cached, err := s.cache.Get(cacheKey)
	if err == nil {
		var books []models.Book
		if err := json.Unmarshal([]byte(cached), &books); err == nil {
			log.Printf("Cache hit for %s", cacheKey)
			return books, nil
		}
	}

	// Cache miss, query database
	books, err := s.repo.GetAll(limit, offset)
	if err != nil {
		log.Printf("Database error in GetAllBooks: %v", err)
		return nil, err
	}

	// Update cache
	serialized, _ := json.Marshal(books)
	if err := s.cache.Set(cacheKey, string(serialized), s.timeout); err != nil {
		log.Printf("Cache set error: %v", err)
	}
	
	return books, nil
}

func (s *BookService) GetBookByID(id uint) (*models.Book, error) {
	cacheKey := fmt.Sprintf("book:%d", id)
	
	// Try cache first
	cached, err := s.cache.Get(cacheKey)
	if err == nil {
		var book models.Book
		if err := json.Unmarshal([]byte(cached), &book); err == nil {
			log.Printf("Cache hit for %s", cacheKey)
			return &book, nil
		}
	}

	// Cache miss, query database
	book, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		log.Printf("Database error in GetBookByID: %v", err)
		return nil, err
	}

	// Update cache
	serialized, _ := json.Marshal(book)
	if err := s.cache.Set(cacheKey, string(serialized), s.timeout); err != nil {
		log.Printf("Cache set error: %v", err)
	}
	
	return book, nil
}

func (s *BookService) CreateBook(book *models.Book) error {
	if err := validateBook(book); err != nil {
		return err
	}

	if err := s.repo.Create(book); err != nil {
		log.Printf("Database error in CreateBook: %v", err)
		return err
	}

	// Invalidate cache
	if err := s.cache.Delete("books:*"); err != nil {
		log.Printf("Cache delete error: %v", err)
	}

	// Publish Kafka event
	event := map[string]interface{}{
		"event_type": "book_created",
		"book":       book,
	}
	eventData, _ := json.Marshal(event)
	if err := kafka.PublishEvent("book_events", eventData); err != nil {
		log.Printf("Failed to publish Kafka event: %v", err)
	}

	return nil
}

// Add similar methods for Update and Delete with proper error handling

// Custom errors
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("book not found")
)

func validateBook(book *models.Book) error {
	if book.Title == "" {
		return fmt.Errorf("%w: title is required", ErrInvalidInput)
	}
	if book.Author == "" {
		return fmt.Errorf("%w: author is required", ErrInvalidInput)
	}
	if book.Year < 0 || book.Year > time.Now().Year()+1 {
		return fmt.Errorf("%w: invalid year", ErrInvalidInput)
	}
	return nil
}