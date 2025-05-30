package line_user

import (
	"net/http"
	lineUserService "opsalert/internal/service/line_user"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *lineUserService.Service
}

func NewHandler(service *lineUserService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) List(c *gin.Context) {
	oaIDStr := c.Query("oa_id")
	if oaIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "oa_id is required"})
		return
	}

	oaID, err := strconv.Atoi(oaIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid oa_id"})
		return
	}

	// ตรวจสอบว่าเป็น admin หรือ staff
	role := c.GetString("role")
	if role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	staffID := c.GetUint("user_id")
	if staffID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	users, err := h.service.GetByOaID(c.Request.Context(), oaID, int(staffID), role)
	if err != nil {
		switch err.Error() {
		case "insufficient permissions to view this OA":
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to view this OA"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// ตรวจสอบว่าเป็น admin หรือ staff
	role := c.GetString("role")
	if role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	staffID := c.GetUint("user_id")
	if staffID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), uint(id), int(staffID), role)
	if err != nil {
		switch err.Error() {
		case "line user not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "line user not found"})
		case "insufficient permissions to view this OA":
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to view this OA"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
