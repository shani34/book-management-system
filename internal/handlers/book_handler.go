package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shani34/book-management-system/internal/models"
	"github.com/shani34/book-management-system/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type BookHandler struct {
	service *services.BookService
	logger  *zap.Logger
}

func NewBookHandler(service *services.BookService, logger *zap.Logger) *BookHandler {
	return &BookHandler{
		service: service,
		logger:  logger.Named("handlers.BookHandler"),
	}
}

// GetBooks godoc
// @Summary List books
// @Description Get paginated list of books
// @Tags books
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.Book
// @Failure 500 {object} map[string]string
// @Router /books [get]
func (h *BookHandler) GetBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	h.logger.Info("Starting GetBooks request",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	books, err := h.service.GetAllBooks(limit, offset)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Warn("No books found", 
				zap.Int("limit", limit),
				zap.Int("offset", offset),
			)
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		} else {
			h.logger.Error("Failed to retrieve books",
				zap.Error(err),
				zap.Int("limit", limit),
				zap.Int("offset", offset),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	h.logger.Info("Successfully retrieved books",
		zap.Int("count", len(books)),
	)
	c.JSON(http.StatusOK, books)
}

// GetBook godoc
// @Summary Get a book
// @Description Get book by ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} models.Book
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid book ID format",
			zap.String("received_id", c.Param("id")),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID format"})
		return
	}

	h.logger.Info("Fetching book", zap.Int("book_id", id))

	book, err := h.service.GetBookByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Warn("Book not found", zap.Int("book_id", id))
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		} else {
			h.logger.Error("Failed to fetch book",
				zap.Int("book_id", id),
				zap.Error(err),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	h.logger.Info("Successfully retrieved book", zap.Int("book_id", id))
	c.JSON(http.StatusOK, book)
}

// CreateBook godoc
// @Summary Create book
// @Description Create new book
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.BookRequest true "Book data"
// @Success 201 {object} models.Book
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var book models.Book
	
	h.logger.Info("Starting CreateBook request")
	
	if err := c.ShouldBindJSON(&book); err != nil {
		h.logger.Warn("Invalid request body",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	h.logger.Debug("Creating book",
		zap.String("title", book.Title),
		zap.String("author", book.Author),
		zap.Int("year", book.Year),
	)

	if err := h.service.CreateBook(&book); err != nil {
		h.logger.Error("Failed to create book",
			zap.Error(err),
			zap.Any("book_data", book),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create book"})
		return
	}

	h.logger.Info("Book created successfully", 
		zap.Uint("book_id", book.ID),
	)
	c.JSON(http.StatusCreated, book)
}

// UpdateBook godoc
// @Summary Update book
// @Description Update existing book
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body models.BookRequest true "Book data"
// @Success 200 {object} models.Book
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid book ID format",
			zap.String("received_id", c.Param("id")),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID format"})
		return
	}

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		h.logger.Warn("Invalid request body for update",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	h.logger.Info("Updating book",
		zap.Int("book_id", id),
		zap.Any("update_data", book),
	)

	if err := h.service.UpdateBook(uint(id), &book); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Warn("Book not found for update",
				zap.Int("book_id", id),
			)
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found, cannot update"})
		} else {
			h.logger.Error("Failed to update book",
				zap.Int("book_id", id),
				zap.Error(err),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update book"})
		}
		return
	}

	h.logger.Info("Book updated successfully", zap.Int("book_id", id))
	c.JSON(http.StatusOK, book)
}

// DeleteBook godoc
// @Summary Delete book
// @Description Delete book by ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("Invalid book ID format",
			zap.String("received_id", c.Param("id")),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID format"})
		return
	}

	h.logger.Info("Deleting book", zap.Int("book_id", id))

	if err := h.service.DeleteBook(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Warn("Book not found for deletion",
				zap.Int("book_id", id),
			)
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		} else {
			h.logger.Error("Failed to delete book",
				zap.Int("book_id", id),
				zap.Error(err),
			)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete book"})
		}
		return
	}

	h.logger.Info("Book deleted successfully", zap.Int("book_id", id))
	c.Status(http.StatusNoContent)
}