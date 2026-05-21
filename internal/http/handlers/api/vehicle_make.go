package api

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"
	"github.com/gin-gonic/gin"
)

type VehicleMakeAPIHandler struct {
	svc vehicle.VehicleMakeService
}

// NewVehicleMakeAPIHandler creates a new VehicleMakeAPIHandler.
func NewVehicleMakeAPIHandler(svc vehicle.VehicleMakeService) *VehicleMakeAPIHandler {
	return &VehicleMakeAPIHandler{svc: svc}
}

func (h *VehicleMakeAPIHandler) ListAll(c *gin.Context) {
	makes, err := h.svc.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"makes": makes})
}

func (h *VehicleMakeAPIHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	makes, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"makes": makes,
		"total": total,
	})
}

func (h *VehicleMakeAPIHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	make, err := h.svc.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if make == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "make not found"})
		return
	}
	c.JSON(http.StatusOK, make)
}

func (h *VehicleMakeAPIHandler) Create(c *gin.Context) {
	var payload struct {
		Name   string `json:"name" binding:"required"`
		Status bool   `json:"status"`
		User   string `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	make, err := h.svc.CreateMake(payload.Name, payload.Status, payload.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, make)
}

func (h *VehicleMakeAPIHandler) Update(c *gin.Context) {
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

	make, err := h.svc.UpdateMake(id, payload.Name, payload.Status, payload.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if make == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "make not found"})
		return
	}
	c.JSON(http.StatusOK, make)
}

func (h *VehicleMakeAPIHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.svc.DeleteMake(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
