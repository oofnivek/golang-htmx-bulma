package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type VehicleFuelHandler struct {
	svc      vehicle.VehicleFuelService
	makeSvc  vehicle.VehicleMakeService
	modelSvc vehicle.VehicleModelService
	ftSvc    vehicle.FuelTypeService
}

func NewVehicleFuelHandler(svc vehicle.VehicleFuelService, makeSvc vehicle.VehicleMakeService, modelSvc vehicle.VehicleModelService, ftSvc vehicle.FuelTypeService) *VehicleFuelHandler {
	return &VehicleFuelHandler{svc: svc, makeSvc: makeSvc, modelSvc: modelSvc, ftSvc: ftSvc}
}

func (h *VehicleFuelHandler) List(c *gin.Context) {
	tz := c.Query("tz")
	if tz == "" {
		tz = "UTC"
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	fuels, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	c.HTML(http.StatusOK, "pages/vehicle_fuels/index.html", gin.H{
		"title":           "Vehicle Fuels",
		"fuels":           fuels,
		"timezone":        tz,
		"timezones":       Timezones,
		"currentPage":     page,
		"pageSize":        pageSize,
		"pageSizeOptions": []int{5, 10, 20, 50},
		"totalPages":      totalPages,
		"total":           total,
		"sortBy":          sortBy,
		"sortOrder":       sortOrder,
	})
}

func (h *VehicleFuelHandler) CreateForm(c *gin.Context) {
	makes, err := h.makeSvc.ListAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	models, err := h.modelSvc.ListAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	fuelTypes, err := h.ftSvc.ListAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "pages/vehicle_fuels/form.html", gin.H{
		"title":         "Add Vehicle Fuel",
		"action":        "/vehicle-fuels",
		"vehicleMakes":  makes,
		"vehicleModels": models,
		"fuelTypes":     fuelTypes,
	})
}

func (h *VehicleFuelHandler) Create(c *gin.Context) {
	vehicleMakeID, _ := strconv.ParseInt(c.PostForm("vehicle_make_id"), 10, 64)
	vehicleModelID, _ := strconv.ParseInt(c.PostForm("vehicle_model_id"), 10, 64)
	fuelTypeID, _ := strconv.ParseInt(c.PostForm("fuel_type_id"), 10, 64)
	fuelTankSize, _ := strconv.ParseFloat(c.PostForm("fuel_tank_size"), 64)
	fuelConsumption, _ := strconv.ParseFloat(c.PostForm("fuel_consumption"), 64)
	status := c.PostForm("status") == "on"

	userEmail := c.GetString("user_email")
	_, err := h.svc.CreateFuel(vehicleMakeID, vehicleModelID, fuelTypeID, fuelTankSize, fuelConsumption, status, userEmail)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicle-fuels")
}

func (h *VehicleFuelHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	f, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if f == nil {
		c.String(http.StatusNotFound, "Vehicle fuel not found")
		return
	}

	makes, err := h.makeSvc.ListAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	models, err := h.modelSvc.ListAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	fuelTypes, err := h.ftSvc.ListAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "pages/vehicle_fuels/form.html", gin.H{
		"title":              "Edit Vehicle Fuel",
		"action":             "/vehicle-fuels/" + idStr,
		"fuel":               f,
		"vehicleMakes":       makes,
		"vehicleModels":      models,
		"fuelTypes":          fuelTypes,
		"currentMakeID":      f.VehicleMakeID,
		"currentModelID":     f.VehicleModelID,
		"currentFuelTypeID":  f.FuelTypeID,
	})
}

func (h *VehicleFuelHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	vehicleMakeID, _ := strconv.ParseInt(c.PostForm("vehicle_make_id"), 10, 64)
	vehicleModelID, _ := strconv.ParseInt(c.PostForm("vehicle_model_id"), 10, 64)
	fuelTypeID, _ := strconv.ParseInt(c.PostForm("fuel_type_id"), 10, 64)
	fuelTankSize, _ := strconv.ParseFloat(c.PostForm("fuel_tank_size"), 64)
	fuelConsumption, _ := strconv.ParseFloat(c.PostForm("fuel_consumption"), 64)
	status := c.PostForm("status") == "on"

	userEmail := c.GetString("user_email")
	_, err := h.svc.UpdateFuel(id, vehicleMakeID, vehicleModelID, fuelTypeID, fuelTankSize, fuelConsumption, status, userEmail)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicle-fuels")
}

func (h *VehicleFuelHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteFuel(id); err != nil {
		slog.Error("failed to delete vehicle fuel", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete vehicle fuel")
		return
	}

	c.Status(http.StatusOK)
}

func (h *VehicleFuelHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tz := c.DefaultQuery("tz", "UTC")

	f, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if f == nil {
		c.String(http.StatusNotFound, "Vehicle fuel not found")
		return
	}

	c.HTML(http.StatusOK, "pages/vehicle_fuels/view.html", gin.H{
		"title": "View Vehicle Fuel",
		"fuel":  f,
		"tz":    tz,
	})
}

func (h *VehicleFuelHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	f, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if f == nil {
		c.String(http.StatusNotFound, "Vehicle fuel not found")
		return
	}

	displayName := fmt.Sprintf("%s / %s / %s", f.VehicleMakeName, f.VehicleModelName, f.FuelTypeName)
	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      displayName,
		"DeleteURL": "/vehicle-fuels/" + idStr,
		"RowID":     "fuel-row-" + idStr,
	})
}
