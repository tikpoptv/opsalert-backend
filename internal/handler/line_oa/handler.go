package line_oa

import (
	"net/http"
	lineOAModel "opsalert/internal/model/line_oa"
	lineOAService "opsalert/internal/service/line_oa"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *lineOAService.Service
}

func NewHandler(service *lineOAService.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var req lineOAModel.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	if err := h.service.Create(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create line official account"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "line official account created successfully"})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req lineOAModel.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, &req); err != nil {
		if err.Error() == "line official account not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "line official account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update line official account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "line official account updated successfully"})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if err.Error() == "line official account not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "line official account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete line official account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "line official account deleted successfully"})
}

func (h *Handler) List(c *gin.Context) {
	oas, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get line official accounts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": oas})
}
