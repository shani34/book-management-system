package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/shani34/book-management-system/internal/models"
	"os"
)

var DB *gorm.DB

func InitDB()(*gorm.DB, error) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto migrate models
	err = DB.AutoMigrate(&models.Book{})
	if err != nil {
		panic("failed to migrate database")
	}

	return DB, nil
}