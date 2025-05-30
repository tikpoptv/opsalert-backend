package main

import (
	"database/sql"
	"fmt"
	"log"
	"opsalert/config"
	"opsalert/internal/handler"
	lineOAHandler "opsalert/internal/handler/line_oa"
	staffHandler "opsalert/internal/handler/staff"
	jwtService "opsalert/internal/jwt"
	lineOARepo "opsalert/internal/repository/line_oa"
	staffRepo "opsalert/internal/repository/staff"
	lineOAService "opsalert/internal/service/line_oa"
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
	lineOARepo := lineOARepo.NewRepository(db)

	jwtSettings := jwtService.DefaultSettings()
	jwtService := jwtService.NewService(jwtSettings)

	staffService := staffService.NewService(staffRepo, jwtService)
	lineOAService := lineOAService.NewService(lineOARepo, config.Get().Domain)

	staffHandler := staffHandler.NewHandler(staffService)
	lineOAHandler := lineOAHandler.NewHandler(lineOAService)

	r := gin.Default()

	handler.SetupRoutes(r, staffHandler, lineOAHandler, jwtService)

	port := config.Get().Port
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
