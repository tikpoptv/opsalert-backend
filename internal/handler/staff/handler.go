package staff

import (
	"net/http"
	staffModel "opsalert/internal/model/staff"
	staffService "opsalert/internal/service/staff"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *staffService.Service
}

func NewHandler(service *staffService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req staffModel.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Register(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "staff registered successfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var req staffModel.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, staff, err := h.service.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"type":  "Bearer",
		"user": gin.H{
			"id":        staff.ID,
			"username":  staff.Username,
			"full_name": staff.FullName,
			"role":      staff.Role,
		},
	})
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	staff, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         staff.ID,
		"username":   staff.Username,
		"full_name":  staff.FullName,
		"role":       staff.Role,
		"is_active":  staff.IsActive,
		"created_at": staff.CreatedAt,
	})
}
