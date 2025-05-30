package main

import (
	"database/sql"
	"fmt"
	"log"
	"opsalert/config"
	"opsalert/internal/handler"
	staffHandler "opsalert/internal/handler/staff"
	jwtService "opsalert/internal/jwt"
	staffRepo "opsalert/internal/repository/staff"
	staffService "opsalert/internal/service/staff"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Failed to load environment: %v", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Get().DBHost,
		config.Get().DBPort,
		config.Get().DBUser,
		config.Get().DBPassword,
		config.Get().DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	staffRepo := staffRepo.NewRepository(db)

	jwtSettings := jwtService.DefaultSettings()
	jwtService := jwtService.NewService(jwtSettings)

	staffService := staffService.NewService(staffRepo, jwtService)

	staffHandler := staffHandler.NewHandler(staffService)

	r := gin.Default()

	handler.SetupRoutes(r, staffHandler, jwtService)

	port := config.Get().Port
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
