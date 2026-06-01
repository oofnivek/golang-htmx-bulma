package api

import (
	"net/http"
	"strconv"
	"time"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type VehicleAPIHandler struct {
	svc vehicle.VehicleService
}

func NewVehicleAPIHandler(svc vehicle.VehicleService) *VehicleAPIHandler {
	return &VehicleAPIHandler{svc: svc}
}

func (h *VehicleAPIHandler) ListAll(c *gin.Context) {
	vehicles, err := h.svc.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"vehicles": vehicles})
}

func (h *VehicleAPIHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	vehicles, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vehicles": vehicles,
		"total":    total,
	})
}

func (h *VehicleAPIHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	v, err := h.svc.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if v == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *VehicleAPIHandler) Create(c *gin.Context) {
	var payload struct {
		VehicleMakeID      int64      `json:"vehicle_make_id" binding:"required"`
		VehicleModelID     int64      `json:"vehicle_model_id" binding:"required"`
		VehicleTypeID      int64      `json:"vehicle_type_id" binding:"required"`
		FuelTypeID         int64      `json:"fuel_type_id" binding:"required"`
		VehicleColorID     int64      `json:"vehicle_color_id" binding:"required"`
		Description        *string    `json:"description"`
		PlateNumber        *string    `json:"plate_number"`
		IUNumber           *string    `json:"iu_number"`
		ChassisNumber      *string    `json:"chassis_number"`
		EngineNumber       *string    `json:"engine_number"`
		NumSeats           int        `json:"num_seats" binding:"required"`
		BootSpace          *string    `json:"boot_space"`
		CarParkID          int64      `json:"car_park_id" binding:"required"`
		AssetOwnerID       int64      `json:"asset_owner_id" binding:"required"`
		VehicleStatusID    int64      `json:"vehicle_status_id" binding:"required"`
		LastServiceDate    *time.Time `json:"last_service_date"`
		LastCleanedDate    *time.Time `json:"last_cleaned_date"`
		LastServiceMileage *int       `json:"last_service_mileage"`
		CurrentMileage     *int       `json:"current_mileage"`
		CurrentFuelLevel   *int       `json:"current_fuel_level"`
		ActiveFrom         *time.Time `json:"active_from"`
		ActiveTo           *time.Time `json:"active_to"`
		User               string     `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v, err := h.svc.CreateVehicle(
		payload.VehicleMakeID, payload.VehicleModelID, payload.VehicleTypeID,
		payload.FuelTypeID, payload.VehicleColorID,
		payload.Description, payload.PlateNumber, payload.IUNumber,
		payload.ChassisNumber, payload.EngineNumber, payload.BootSpace,
		payload.NumSeats, payload.CarParkID, payload.AssetOwnerID, payload.VehicleStatusID,
		payload.LastServiceDate, payload.LastCleanedDate, payload.ActiveFrom, payload.ActiveTo,
		payload.LastServiceMileage, payload.CurrentMileage, payload.CurrentFuelLevel,
		payload.User,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *VehicleAPIHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var payload struct {
		VehicleMakeID      int64      `json:"vehicle_make_id" binding:"required"`
		VehicleModelID     int64      `json:"vehicle_model_id" binding:"required"`
		VehicleTypeID      int64      `json:"vehicle_type_id" binding:"required"`
		FuelTypeID         int64      `json:"fuel_type_id" binding:"required"`
		VehicleColorID     int64      `json:"vehicle_color_id" binding:"required"`
		Description        *string    `json:"description"`
		PlateNumber        *string    `json:"plate_number"`
		IUNumber           *string    `json:"iu_number"`
		ChassisNumber      *string    `json:"chassis_number"`
		EngineNumber       *string    `json:"engine_number"`
		NumSeats           int        `json:"num_seats" binding:"required"`
		BootSpace          *string    `json:"boot_space"`
		CarParkID          int64      `json:"car_park_id" binding:"required"`
		AssetOwnerID       int64      `json:"asset_owner_id" binding:"required"`
		VehicleStatusID    int64      `json:"vehicle_status_id" binding:"required"`
		LastServiceDate    *time.Time `json:"last_service_date"`
		LastCleanedDate    *time.Time `json:"last_cleaned_date"`
		LastServiceMileage *int       `json:"last_service_mileage"`
		CurrentMileage     *int       `json:"current_mileage"`
		CurrentFuelLevel   *int       `json:"current_fuel_level"`
		ActiveFrom         *time.Time `json:"active_from"`
		ActiveTo           *time.Time `json:"active_to"`
		User               string     `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v, err := h.svc.UpdateVehicle(
		id,
		payload.VehicleMakeID, payload.VehicleModelID, payload.VehicleTypeID,
		payload.FuelTypeID, payload.VehicleColorID,
		payload.Description, payload.PlateNumber, payload.IUNumber,
		payload.ChassisNumber, payload.EngineNumber, payload.BootSpace,
		payload.NumSeats, payload.CarParkID, payload.AssetOwnerID, payload.VehicleStatusID,
		payload.LastServiceDate, payload.LastCleanedDate, payload.ActiveFrom, payload.ActiveTo,
		payload.LastServiceMileage, payload.CurrentMileage, payload.CurrentFuelLevel,
		payload.User,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if v == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *VehicleAPIHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.svc.DeleteVehicle(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
