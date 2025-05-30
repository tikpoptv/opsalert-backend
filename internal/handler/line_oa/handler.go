package line_oa

import (
	"net/http"
	lineOAModel "opsalert/internal/model/line_oa"
	lineOAService "opsalert/internal/service/line_oa"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *lineOAService.Service
}

func NewHandler(service *lineOAService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateOA(c *gin.Context) {
	var req lineOAModel.CreateLineOARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateOA(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "line official account created successfully"})
}
