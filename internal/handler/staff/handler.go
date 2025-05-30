package staff

import (
	"net/http"
	staffModel "opsalert/internal/model/staff"
	staffService "opsalert/internal/service/staff"
	"strconv"
	"strings"

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

func (h *Handler) GetAccounts(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can access this endpoint"})
		return
	}

	accounts, err := h.service.GetAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}

func (h *Handler) GetAccountByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	staff, err := h.service.GetProfile(uint(id))
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

func (h *Handler) UpdateStaff(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	var req staffModel.UpdateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateStaff(uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "staff updated successfully"})
}

func (h *Handler) SetPermissions(c *gin.Context) {
	var req staffModel.PermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}

	if err := h.service.SetPermissions(c.Request.Context(), &req); err != nil {
		switch {
		case err.Error() == "staff not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "staff not found"})
		case err.Error() == "cannot set permissions for admin":
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot set permissions for admin"})
		case strings.Contains(err.Error(), "OA with ID"):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "staff permissions updated successfully"})
}

func (h *Handler) GetStaffPermissions(c *gin.Context) {
	staffIDStr := c.Param("staff_id")
	staffID, err := strconv.Atoi(staffIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff id"})
		return
	}

	permissions, err := h.service.GetStaffPermissions(c.Request.Context(), staffID)
	if err != nil {
		switch err.Error() {
		case "staff not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "staff not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": permissions})
}

func (h *Handler) DeleteStaffPermissions(c *gin.Context) {
	staffIDStr := c.Param("id")
	staffID, err := strconv.Atoi(staffIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid staff id"})
		return
	}

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

	if err := h.service.DeleteStaffPermissions(c.Request.Context(), staffID, oaID); err != nil {
		switch err.Error() {
		case "staff not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "staff not found"})
		case "OA not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "OA not found"})
		case "staff does not have permission for this OA":
			c.JSON(http.StatusBadRequest, gin.H{"error": "staff does not have permission for this OA"})
		case "cannot delete permissions for admin":
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete permissions for admin"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "staff permission deleted successfully"})
}
