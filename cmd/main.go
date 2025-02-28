package main

import (
    "github.com/shani34/book-management-system/api"
    "github.com/shani34/book-management-system/config"
    "log"
)

// @title           Book Management API
// @version         1.0
// @description     This is a REST API for managing books
// @host            localhost:8080
// @BasePath        /api/v1
// @schemes         http
func main() {
    // Load environment variables
    config.LoadEnv()
    
    
    // Create router with middleware
    router := api.SetupRouter()
    
    // Start server
    log.Println("Server starting on :8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}