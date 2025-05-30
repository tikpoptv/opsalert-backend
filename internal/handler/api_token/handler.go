package api_token

import (
	"net/http"
	apiTokenService "opsalert/internal/service/api_token"
	"strconv"

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

	// Create API token
	token, err := h.service.Create(c.Request.Context(), userID.(int), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusCreated, token)
}

func (h *Handler) Reset(c *gin.Context) {
	tokenID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token id"})
		return
	}

	// Get user ID and role from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Reset token
	token, err := h.service.ResetToken(c.Request.Context(), tokenID, userID.(int), role.(string))
	if err != nil {
		if err == apiTokenService.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to reset this token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset token"})
		return
	}

	c.JSON(http.StatusOK, token)
}

type UpdateStatusRequest struct {
	IsActive bool `json:"is_active" binding:"required"`
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	tokenID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token id"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID and role from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Update token status
	token, err := h.service.UpdateStatus(c.Request.Context(), tokenID, userID.(int), role.(string), req.IsActive)
	if err != nil {
		if err == apiTokenService.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to update this token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update token status"})
		return
	}

	c.JSON(http.StatusOK, token)
}
