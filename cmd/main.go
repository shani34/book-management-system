package main

import (
    "github.com/shani34/book-management-system/api"
    "github.com/shani34/book-management-system/config"
    "github.com/shani34/book-management-system/pkg/db"
    "github.com/shani34/book-management-system/pkg/kafka"
    "github.com/shani34/book-management-system/pkg/redis"
    "log"
)

// @title Book Management API
// @version 1.0
// @description This is a REST API for managing books
// @host localhost:8080
// @BasePath /api/v1
func main() {
    // Load environment variables
    config.LoadEnv()
    
    // Initialize dependencies
    db.InitDB()
    redis.InitRedis()
    kafka.InitKafkaProducer()
    
    // Create router with middleware
    router := api.SetupRouter()
    
    // Start server
    log.Println("Server starting on :8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}