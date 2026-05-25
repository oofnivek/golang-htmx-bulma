package api

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type VehicleStatusAPIHandler struct {
	svc vehicle.VehicleStatusService
}

func NewVehicleStatusAPIHandler(svc vehicle.VehicleStatusService) *VehicleStatusAPIHandler {
	return &VehicleStatusAPIHandler{svc: svc}
}

func (h *VehicleStatusAPIHandler) ListAll(c *gin.Context) {
	statuses, err := h.svc.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"statuses": statuses})
}

func (h *VehicleStatusAPIHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	statuses, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statuses": statuses,
		"total":    total,
	})
}

func (h *VehicleStatusAPIHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	vs, err := h.svc.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if vs == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "status not found"})
		return
	}
	c.JSON(http.StatusOK, vs)
}

func (h *VehicleStatusAPIHandler) Create(c *gin.Context) {
	var payload struct {
		Substatus string `json:"substatus" binding:"required"`
		IsActive  bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vs, err := h.svc.CreateStatus(payload.Substatus, payload.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, vs)
}

func (h *VehicleStatusAPIHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var payload struct {
		Substatus string `json:"substatus" binding:"required"`
		IsActive  bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vs, err := h.svc.UpdateStatus(id, payload.Substatus, payload.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if vs == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "status not found"})
		return
	}
	c.JSON(http.StatusOK, vs)
}

func (h *VehicleStatusAPIHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.svc.DeleteStatus(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
