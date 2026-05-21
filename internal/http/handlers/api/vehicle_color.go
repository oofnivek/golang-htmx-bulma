package api

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"
	"github.com/gin-gonic/gin"
)

type VehicleColorAPIHandler struct {
	svc vehicle.VehicleColorService
}

// NewVehicleColorAPIHandler creates a new VehicleColorAPIHandler.
func NewVehicleColorAPIHandler(svc vehicle.VehicleColorService) *VehicleColorAPIHandler {
	return &VehicleColorAPIHandler{svc: svc}
}

func (h *VehicleColorAPIHandler) ListAll(c *gin.Context) {
	colors, err := h.svc.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"colors": colors})
}

func (h *VehicleColorAPIHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	colors, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"colors": colors,
		"total":  total,
	})
}

func (h *VehicleColorAPIHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	color, err := h.svc.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if color == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "color not found"})
		return
	}
	c.JSON(http.StatusOK, color)
}

func (h *VehicleColorAPIHandler) Create(c *gin.Context) {
	var payload struct {
		Name   string `json:"name" binding:"required"`
		Status bool   `json:"status"`
		User   string `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	color, err := h.svc.CreateColor(payload.Name, payload.Status, payload.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, color)
}

func (h *VehicleColorAPIHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var payload struct {
		Name   string `json:"name" binding:"required"`
		Status bool   `json:"status"`
		User   string `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	color, err := h.svc.UpdateColor(id, payload.Name, payload.Status, payload.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if color == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "color not found"})
		return
	}
	c.JSON(http.StatusOK, color)
}

func (h *VehicleColorAPIHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.svc.DeleteColor(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
