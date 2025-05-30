package handler

import (
	staffHandler "opsalert/internal/handler/staff"
	jwtService "opsalert/internal/jwt"
	"opsalert/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, staffHandler *staffHandler.Handler, jwtService *jwtService.Service) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	v1 := r.Group("/api/v1")
	{
		staff := v1.Group("/staff")
		{
			// Public routes
			staff.POST("/login", staffHandler.Login)

			// Protected routes
			staff.Use(middleware.AuthMiddleware(jwtService))
			staff.GET("/me", staffHandler.GetProfile)
			staff.GET("/accounts", middleware.AdminOnly(), staffHandler.GetAccounts)
			staff.GET("/accounts/:id", middleware.AdminOnly(), staffHandler.GetAccountByID)
			staff.PUT("/accounts/:id", middleware.AdminOnly(), staffHandler.UpdateStaff)
			staff.POST("/register", middleware.AdminOnly(), staffHandler.Register)
		}

		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}
}
