package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shani34/book-management-system/internal/handlers"
	"github.com/shani34/book-management-system/internal/middleware"
	"github.com/shani34/book-management-system/internal/repositories"
	"github.com/shani34/book-management-system/internal/services"
	"github.com/shani34/book-management-system/pkg/db"
	"github.com/shani34/book-management-system/pkg/kafka"
	"github.com/shani34/book-management-system/pkg/redis"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()

    // Middleware
    router.Use(middleware.Logger())
    router.Use(gin.Recovery())
    
    // Initialize dependencies
	db, err:=db.InitDB()
	if err!=nil{
		gin.Logger()
	}
	redisClient, err := redis.InitRedis()
	if err!=nil{
		gin.Logger()
	}
	kafka.InitKafkaProducer()
    bookRepo := repositories.NewBookRepository(db)
    
    bookService := services.NewBookService(bookRepo, redisClient)
    bookHandler := handlers.NewBookHandler(bookService)
    
    // API routes
    v1 := router.Group("/api/v1")
    {
        books := v1.Group("/books")
        {
            books.GET("", bookHandler.GetBooks)
            books.POST("", bookHandler.CreateBook)
            books.GET("/:id", bookHandler.GetBook)
            books.PUT("/:id", bookHandler.UpdateBook)
            books.DELETE("/:id", bookHandler.DeleteBook)
        }
    }
    
    // Swagger documentation
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    return router
}