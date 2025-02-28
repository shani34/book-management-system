package models

import (
	"time"
)

type Book struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Author    string    `gorm:"not null" json:"author"`
	Year      int       `json:"year"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Swagger model documentation
type BookRequest struct {
	Title  string `json:"title" example:"The Go Programming Language"`
	Author string `json:"author" example:"Alan A. A. Donovan"`
	Year   int    `json:"year" example:"2015"`
}