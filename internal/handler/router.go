package handler

import (
	apiTokenHandler "opsalert/internal/handler/api_token"
	lineOAHandler "opsalert/internal/handler/line_oa"
	lineUserHandler "opsalert/internal/handler/line_user"
	staffHandler "opsalert/internal/handler/staff"
	webhookHandler "opsalert/internal/handler/webhook"
	jwtService "opsalert/internal/jwt"
	"opsalert/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, staffHandler *staffHandler.Handler, lineOAHandler *lineOAHandler.Handler, lineUserHandler *lineUserHandler.Handler, apiTokenHandler *apiTokenHandler.Handler, webhookHandler *webhookHandler.Handler, jwtService *jwtService.Service) {
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
			staff.POST("/permissions", middleware.AdminOnly(), staffHandler.SetPermissions)
			staff.GET("/permissions/:staff_id", middleware.AdminOnly(), staffHandler.GetStaffPermissions)
			staff.DELETE("/permissions/:id", middleware.AdminOnly(), staffHandler.DeleteStaffPermissions)
		}

		oa := v1.Group("/oa")
		{
			oa.Use(middleware.AuthMiddleware(jwtService))
			oa.POST("", middleware.AdminOnly(), lineOAHandler.Create)
			oa.PUT("/:id", middleware.StaffOnly(), lineOAHandler.Update)
			oa.DELETE("/:id", middleware.AdminOnly(), lineOAHandler.Delete)
			oa.GET("", lineOAHandler.List)
		}

		lineUsers := v1.Group("/line-users")
		{
			lineUsers.Use(middleware.AuthMiddleware(jwtService))
			lineUsers.GET("", lineUserHandler.List)
			lineUsers.GET("/:id", lineUserHandler.GetByID)
		}

		// API Tokens
		apiTokens := v1.Group("/api-tokens")
		apiTokens.Use(middleware.AuthMiddleware(jwtService))
		{
			apiTokens.POST("", apiTokenHandler.Create)
			apiTokens.POST("/:id/reset", apiTokenHandler.Reset)
			apiTokens.PUT("/:id/status", apiTokenHandler.UpdateStatus)
		}

		// Webhook routes
		webhooks := v1.Group("/webhooks")
		{
			webhooks.POST("/line/:oa_id", webhookHandler.HandleLineWebhook)
		}

		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}
}
