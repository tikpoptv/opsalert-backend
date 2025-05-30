package handler

import (
	staffHandler "opsalert/internal/handler/staff"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, staffHandler *staffHandler.Handler) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	v1 := r.Group("/api/v1")
	{
		staff := v1.Group("/staff")
		{
			staff.POST("/register", staffHandler.Register)
			staff.POST("/login", staffHandler.Login)
		}

		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}
}
