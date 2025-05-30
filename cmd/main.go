package main

import (
	"database/sql"
	"fmt"
	"log"
	"opsalert/config"
	"opsalert/internal/handler"
	lineOAHandler "opsalert/internal/handler/line_oa"
	lineUserHandler "opsalert/internal/handler/line_user"
	staffHandler "opsalert/internal/handler/staff"
	jwtService "opsalert/internal/jwt"
	lineOARepo "opsalert/internal/repository/line_oa"
	lineUserRepo "opsalert/internal/repository/line_user"
	staffRepo "opsalert/internal/repository/staff"
	lineOAService "opsalert/internal/service/line_oa"
	lineUserService "opsalert/internal/service/line_user"
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
	lineUserRepo := lineUserRepo.NewRepository(db)

	jwtSettings := jwtService.DefaultSettings()
	jwtService := jwtService.NewService(jwtSettings)

	staffService := staffService.NewService(staffRepo, jwtService)
	lineOAService := lineOAService.NewService(lineOARepo, config.Get().Domain)
	lineUserService := lineUserService.NewService(lineUserRepo, staffRepo)

	staffHandler := staffHandler.NewHandler(staffService)
	lineOAHandler := lineOAHandler.NewHandler(lineOAService)
	lineUserHandler := lineUserHandler.NewHandler(lineUserService)

	r := gin.Default()

	handler.SetupRoutes(r, staffHandler, lineOAHandler, lineUserHandler, nil, nil, jwtService)

	port := config.Get().Port
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
