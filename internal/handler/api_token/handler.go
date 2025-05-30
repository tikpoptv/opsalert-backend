package api_token

import (
	"net/http"
	apiTokenService "opsalert/internal/service/api_token"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *apiTokenService.Service
}

func NewHandler(service *apiTokenService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

type CreateRequest struct {
	OAID int    `json:"oa_id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Check if user has permission for this OA
	hasPermission, err := h.service.CheckStaffOAPermission(c.Request.Context(), userID.(int), req.OAID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check permission"})
		return
	}
	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	// Create API token
	token, err := h.service.Create(c.Request.Context(), userID.(int), req.OAID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusCreated, token)
}
