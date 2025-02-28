package services

import (
	"encoding/json"
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
	cache   *redis.RedisClient
	timeout time.Duration
}

func NewBookService(repo *repositories.BookRepository, cache *redis.RedisClient) *BookService {
	return &BookService{
		repo:    repo,
		cache:   cache,
		timeout: 10 * time.Minute,
	}
}

func (s *BookService) GetAllBooks(limit, offset int) ([]models.Book, error) {
	cacheKey := fmt.Sprintf("books:%d:%d", limit, offset)
	
	if cached, err := s.cache.Get(cacheKey); err == nil {
		var books []models.Book
		if json.Unmarshal([]byte(cached), &books) == nil {
			return books, nil
		}
	}

	books, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	if serialized, err := json.Marshal(books); err == nil {
		s.cache.Set(cacheKey, string(serialized), s.timeout)
	}
	
	return books, nil
}

func (s *BookService) GetBookByID(id uint) (*models.Book, error) {
	cacheKey := fmt.Sprintf("book:%d", id)
	
	if cached, err := s.cache.Get(cacheKey); err == nil {
		var book models.Book
		if json.Unmarshal([]byte(cached), &book) == nil {
			return &book, nil
		}
	}

	book, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if serialized, err := json.Marshal(book); err == nil {
		s.cache.Set(cacheKey, string(serialized), s.timeout)
	}
	
	return book, nil
}

func (s *BookService) CreateBook(book *models.Book) error {
	if err := validateBook(book); err != nil {
		return err
	}

	if err := s.repo.Create(book); err != nil {
		return err
	}

	s.cache.Delete("books:*")
	s.publishKafkaEvent("book_created", book)
	return nil
}

func (s *BookService) UpdateBook(id uint, book *models.Book) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := validateBook(book); err != nil {
		return err
	}

	book.ID = id
	book.CreatedAt = existing.CreatedAt
	
	if err := s.repo.Update(book); err != nil {
		return err
	}

	s.cache.Delete(fmt.Sprintf("book:%d", id), "books:*")
	s.publishKafkaEvent("book_updated", book)
	return nil
}

func (s *BookService) DeleteBook(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	s.cache.Delete(fmt.Sprintf("book:%d", id), "books:*")
	s.publishKafkaEvent("book_deleted", map[string]interface{}{"id": id})
	return nil
}

func validateBook(book *models.Book) error {
	if book.Title == "" {
		return fmt.Errorf("title is required")
	}
	if book.Author == "" {
		return fmt.Errorf("author is required")
	}
	if book.Year < 0 || book.Year > time.Now().Year()+1 {
		return fmt.Errorf("invalid year")
	}
	return nil
}

func (s *BookService) publishKafkaEvent(eventType string, payload interface{}) {
	event := map[string]interface{}{
		"event_type": eventType,
		"payload":    payload,
		"timestamp":  time.Now().UTC(),
	}

	if eventData, err := json.Marshal(event); err == nil {
		if err := kafka.PublishEvent("book_events", eventData); err != nil {
			log.Printf("Failed to publish Kafka event: %v", err)
		}
	}
}