package main

import (
    "github.com/shani34/book-management-system/api"
    "github.com/shani34/book-management-system/config"
	_ "github.com/shani34/book-management-system/docs"
    "log"
)

// @title Book Management API
// @version 1.0
// @description REST API for managing books with Redis caching and Kafka integration
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@bookapi.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host https://book-management-system-production-7d0e.up.railway.app/
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @swagger 2.0  // <-- Add this line to specify Swagger version
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