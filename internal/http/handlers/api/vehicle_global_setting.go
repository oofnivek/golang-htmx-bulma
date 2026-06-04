package api

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type VehicleGlobalSettingAPIHandler struct {
	svc vehicle.VehicleGlobalSettingService
}

func NewVehicleGlobalSettingAPIHandler(svc vehicle.VehicleGlobalSettingService) *VehicleGlobalSettingAPIHandler {
	return &VehicleGlobalSettingAPIHandler{svc: svc}
}

func (h *VehicleGlobalSettingAPIHandler) ListAll(c *gin.Context) {
	settings, err := h.svc.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

func (h *VehicleGlobalSettingAPIHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	settings, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"settings": settings,
		"total":    total,
	})
}

func (h *VehicleGlobalSettingAPIHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	setting, err := h.svc.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if setting == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "setting not found"})
		return
	}
	c.JSON(http.StatusOK, setting)
}

func (h *VehicleGlobalSettingAPIHandler) Create(c *gin.Context) {
	var payload struct {
		Key         string  `json:"key"   binding:"required"`
		Value       string  `json:"value" binding:"required"`
		Remark      *string `json:"remark"`
		CountryCode *string `json:"country_code"`
		User        string  `json:"user"  binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setting, err := h.svc.CreateVehicleGlobalSetting(payload.Key, payload.Value, payload.Remark, payload.CountryCode, payload.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, setting)
}

func (h *VehicleGlobalSettingAPIHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var payload struct {
		Key         string  `json:"key"   binding:"required"`
		Value       string  `json:"value" binding:"required"`
		Remark      *string `json:"remark"`
		CountryCode *string `json:"country_code"`
		User        string  `json:"user"  binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setting, err := h.svc.UpdateVehicleGlobalSetting(id, payload.Key, payload.Value, payload.Remark, payload.CountryCode, payload.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if setting == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "setting not found"})
		return
	}
	c.JSON(http.StatusOK, setting)
}

func (h *VehicleGlobalSettingAPIHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.svc.DeleteVehicleGlobalSetting(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
