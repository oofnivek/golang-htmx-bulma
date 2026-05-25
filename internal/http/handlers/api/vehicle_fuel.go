package api

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type VehicleFuelAPIHandler struct {
	svc vehicle.VehicleFuelService
}

func NewVehicleFuelAPIHandler(svc vehicle.VehicleFuelService) *VehicleFuelAPIHandler {
	return &VehicleFuelAPIHandler{svc: svc}
}

func (h *VehicleFuelAPIHandler) ListAll(c *gin.Context) {
	fuels, err := h.svc.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"fuels": fuels})
}

func (h *VehicleFuelAPIHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	fuels, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"fuels": fuels, "total": total})
}

func (h *VehicleFuelAPIHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	f, err := h.svc.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if f == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vehicle fuel not found"})
		return
	}
	c.JSON(http.StatusOK, f)
}

func (h *VehicleFuelAPIHandler) Create(c *gin.Context) {
	var payload struct {
		VehicleMakeID   int64   `json:"vehicle_make_id"   binding:"required"`
		VehicleModelID  int64   `json:"vehicle_model_id"  binding:"required"`
		FuelTypeID      int64   `json:"fuel_type_id"      binding:"required"`
		FuelTankSize    float64 `json:"fuel_tank_size"    binding:"required"`
		FuelConsumption float64 `json:"fuel_consumption"  binding:"required"`
		Status          bool    `json:"status"`
		User            string  `json:"user"              binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	f, err := h.svc.CreateFuel(
		payload.VehicleMakeID, payload.VehicleModelID, payload.FuelTypeID,
		payload.FuelTankSize, payload.FuelConsumption,
		payload.Status, payload.User,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, f)
}

func (h *VehicleFuelAPIHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var payload struct {
		VehicleMakeID   int64   `json:"vehicle_make_id"   binding:"required"`
		VehicleModelID  int64   `json:"vehicle_model_id"  binding:"required"`
		FuelTypeID      int64   `json:"fuel_type_id"      binding:"required"`
		FuelTankSize    float64 `json:"fuel_tank_size"    binding:"required"`
		FuelConsumption float64 `json:"fuel_consumption"  binding:"required"`
		Status          bool    `json:"status"`
		User            string  `json:"user"              binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	f, err := h.svc.UpdateFuel(
		id,
		payload.VehicleMakeID, payload.VehicleModelID, payload.FuelTypeID,
		payload.FuelTankSize, payload.FuelConsumption,
		payload.Status, payload.User,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if f == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vehicle fuel not found"})
		return
	}
	c.JSON(http.StatusOK, f)
}

func (h *VehicleFuelAPIHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.svc.DeleteFuel(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
