package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/shani34/book-management-system/internal/models"
	"github.com/shani34/book-management-system/internal/services"
	"net/http"
	"strconv"
)

type BookHandler struct {
	service *services.BookService
}

func NewBookHandler(service *services.BookService) *BookHandler {
	return &BookHandler{service: service}
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

	books, err := h.service.GetAllBooks(limit, offset)
	if handleError(c, err) {
		return
	}

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
	if handleError(c, err) {
		return
	}

	book, err := h.service.GetBookByID(uint(id))
	if handleError(c, err) {
		return
	}

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
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.CreateBook(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	if handleError(c, err) {
		return
	}

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdateBook(uint(id), &book); err != nil {
		handleError(c, err)
		return
	}

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
	if handleError(c, err) {
		return
	}

	if err := h.service.DeleteBook(uint(id)); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func handleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

    c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	
	 return true
}