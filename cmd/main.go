package main

import (
	"fmt"
	"log"
	"opsalert/config"
	"opsalert/internal/db"
	"opsalert/internal/handler"
	"opsalert/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// Set gin mode
	if config.AppConfig.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	r := gin.Default()

	// Apply middleware
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// Register routes
	handler.RegisterRoutes(r)

	// Start server
	port := fmt.Sprintf(":%s", config.AppConfig.Port)
	log.Printf("Server starting on port %s", config.AppConfig.Port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
