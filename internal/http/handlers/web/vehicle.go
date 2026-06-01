package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	svc       vehicle.VehicleService
	makeSvc   vehicle.VehicleMakeService
	modelSvc  vehicle.VehicleModelService
	typeSvc   vehicle.VehicleTypeService
	fuelSvc   vehicle.FuelTypeService
	colorSvc  vehicle.VehicleColorService
	parkSvc   vehicle.CarParkService
	ownerSvc  vehicle.CarAssetOwnerService
	statusSvc vehicle.VehicleStatusService
}

func NewVehicleHandler(
	svc vehicle.VehicleService,
	makeSvc vehicle.VehicleMakeService,
	modelSvc vehicle.VehicleModelService,
	typeSvc vehicle.VehicleTypeService,
	fuelSvc vehicle.FuelTypeService,
	colorSvc vehicle.VehicleColorService,
	parkSvc vehicle.CarParkService,
	ownerSvc vehicle.CarAssetOwnerService,
	statusSvc vehicle.VehicleStatusService,
) *VehicleHandler {
	return &VehicleHandler{
		svc:       svc,
		makeSvc:   makeSvc,
		modelSvc:  modelSvc,
		typeSvc:   typeSvc,
		fuelSvc:   fuelSvc,
		colorSvc:  colorSvc,
		parkSvc:   parkSvc,
		ownerSvc:  ownerSvc,
		statusSvc: statusSvc,
	}
}

func (h *VehicleHandler) List(c *gin.Context) {
	tz := c.Query("tz")
	if tz == "" {
		tz = "UTC"
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	vehicles, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		slog.Error("failed to list vehicles", "error", err)
		c.String(http.StatusInternalServerError, "Failed to list vehicles")
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	start := (page-1)*pageSize + 1
	if total == 0 {
		start = 0
	}
	end := page * pageSize
	if end > total {
		end = total
	}

	c.HTML(http.StatusOK, "pages/vehicles/index.html", gin.H{
		"title":           "Vehicles",
		"vehicles":        vehicles,
		"timezone":        tz,
		"timezones":       Timezones,
		"currentPage":     page,
		"pageSize":        pageSize,
		"pageSizeOptions": []int{5, 10, 20, 50},
		"totalPages":      totalPages,
		"total":           total,
		"start":           start,
		"end":             end,
		"sortBy":          sortBy,
		"sortOrder":       sortOrder,
		"pageWindow":      paginationWindow(page, totalPages),
	})
}

func (h *VehicleHandler) loadDropdowns(c *gin.Context) (gin.H, error) {
	makes, err := h.makeSvc.ListAll()
	if err != nil {
		slog.Error("failed to list vehicle makes for form", "error", err)
		return nil, err
	}
	models, err := h.modelSvc.ListAll()
	if err != nil {
		slog.Error("failed to list vehicle models for form", "error", err)
		return nil, err
	}
	types, err := h.typeSvc.ListAll()
	if err != nil {
		slog.Error("failed to list vehicle types for form", "error", err)
		return nil, err
	}
	fuels, err := h.fuelSvc.ListAll()
	if err != nil {
		slog.Error("failed to list fuel types for form", "error", err)
		return nil, err
	}
	colors, err := h.colorSvc.ListAll()
	if err != nil {
		slog.Error("failed to list vehicle colors for form", "error", err)
		return nil, err
	}
	parks, err := h.parkSvc.ListAll()
	if err != nil {
		slog.Error("failed to list car parks for form", "error", err)
		return nil, err
	}
	owners, err := h.ownerSvc.ListAll()
	if err != nil {
		slog.Error("failed to list asset owners for form", "error", err)
		return nil, err
	}
	statuses, err := h.statusSvc.ListAll()
	if err != nil {
		slog.Error("failed to list vehicle statuses for form", "error", err)
		return nil, err
	}
	return gin.H{
		"makes":    makes,
		"models":   models,
		"types":    types,
		"fuels":    fuels,
		"colors":   colors,
		"parks":    parks,
		"owners":   owners,
		"statuses": statuses,
	}, nil
}

func (h *VehicleHandler) CreateForm(c *gin.Context) {
	dropdowns, err := h.loadDropdowns(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}
	dropdowns["title"] = "Add Vehicle"
	dropdowns["action"] = "/vehicles"
	c.HTML(http.StatusOK, "pages/vehicles/form.html", dropdowns)
}

func (h *VehicleHandler) Create(c *gin.Context) {
	userEmail := c.GetString("user_email")

	vehicleMakeID, _ := strconv.ParseInt(c.PostForm("vehicle_make_id"), 10, 64)
	vehicleModelID, _ := strconv.ParseInt(c.PostForm("vehicle_model_id"), 10, 64)
	vehicleTypeID, _ := strconv.ParseInt(c.PostForm("vehicle_type_id"), 10, 64)
	fuelTypeID, _ := strconv.ParseInt(c.PostForm("fuel_type_id"), 10, 64)
	vehicleColorID, _ := strconv.ParseInt(c.PostForm("vehicle_color_id"), 10, 64)
	carParkID, _ := strconv.ParseInt(c.PostForm("car_park_id"), 10, 64)
	assetOwnerID, _ := strconv.ParseInt(c.PostForm("asset_owner_id"), 10, 64)
	vehicleStatusID, _ := strconv.ParseInt(c.PostForm("vehicle_status_id"), 10, 64)
	numSeats, _ := strconv.Atoi(c.PostForm("num_seats"))

	description := nullableString(c.PostForm("description"))
	plateNumber := nullableString(c.PostForm("plate_number"))
	iuNumber := nullableString(c.PostForm("iu_number"))
	chassisNumber := nullableString(c.PostForm("chassis_number"))
	engineNumber := nullableString(c.PostForm("engine_number"))
	bootSpace := nullableString(c.PostForm("boot_space"))

	lastServiceDate := parseDateTimeLocal(c.PostForm("last_service_date"))
	lastCleanedDate := parseDateTimeLocal(c.PostForm("last_cleaned_date"))
	activeFrom := parseDateTimeLocal(c.PostForm("active_from"))
	activeTo := parseDateTimeLocal(c.PostForm("active_to"))

	lastServiceMileage := nullableInt(c.PostForm("last_service_mileage"))
	currentMileage := nullableInt(c.PostForm("current_mileage"))
	currentFuelLevel := nullableInt(c.PostForm("current_fuel_level"))

	_, err := h.svc.CreateVehicle(
		vehicleMakeID, vehicleModelID, vehicleTypeID, fuelTypeID, vehicleColorID,
		description, plateNumber, iuNumber, chassisNumber, engineNumber, bootSpace,
		numSeats, carParkID, assetOwnerID, vehicleStatusID,
		lastServiceDate, lastCleanedDate, activeFrom, activeTo,
		lastServiceMileage, currentMileage, currentFuelLevel,
		userEmail,
	)
	if err != nil {
		slog.Error("failed to create vehicle", "error", err)
		c.String(http.StatusInternalServerError, "Failed to create vehicle")
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicles")
}

func (h *VehicleHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	v, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch vehicle", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch vehicle")
		return
	}
	if v == nil {
		c.String(http.StatusNotFound, "Vehicle not found")
		return
	}

	dropdowns, err := h.loadDropdowns(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}
	dropdowns["title"] = "Edit Vehicle"
	dropdowns["action"] = "/vehicles/" + idStr
	dropdowns["vehicle"] = v
	dropdowns["lastServiceDateStr"] = formatDateTimeLocal(v.LastServiceDate)
	dropdowns["lastCleanedDateStr"] = formatDateTimeLocal(v.LastCleanedDate)
	dropdowns["activeFromStr"] = formatDateTimeLocal(v.ActiveFrom)
	dropdowns["activeToStr"] = formatDateTimeLocal(v.ActiveTo)

	c.HTML(http.StatusOK, "pages/vehicles/form.html", dropdowns)
}

func (h *VehicleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	userEmail := c.GetString("user_email")

	vehicleMakeID, _ := strconv.ParseInt(c.PostForm("vehicle_make_id"), 10, 64)
	vehicleModelID, _ := strconv.ParseInt(c.PostForm("vehicle_model_id"), 10, 64)
	vehicleTypeID, _ := strconv.ParseInt(c.PostForm("vehicle_type_id"), 10, 64)
	fuelTypeID, _ := strconv.ParseInt(c.PostForm("fuel_type_id"), 10, 64)
	vehicleColorID, _ := strconv.ParseInt(c.PostForm("vehicle_color_id"), 10, 64)
	carParkID, _ := strconv.ParseInt(c.PostForm("car_park_id"), 10, 64)
	assetOwnerID, _ := strconv.ParseInt(c.PostForm("asset_owner_id"), 10, 64)
	vehicleStatusID, _ := strconv.ParseInt(c.PostForm("vehicle_status_id"), 10, 64)
	numSeats, _ := strconv.Atoi(c.PostForm("num_seats"))

	description := nullableString(c.PostForm("description"))
	plateNumber := nullableString(c.PostForm("plate_number"))
	iuNumber := nullableString(c.PostForm("iu_number"))
	chassisNumber := nullableString(c.PostForm("chassis_number"))
	engineNumber := nullableString(c.PostForm("engine_number"))
	bootSpace := nullableString(c.PostForm("boot_space"))

	lastServiceDate := parseDateTimeLocal(c.PostForm("last_service_date"))
	lastCleanedDate := parseDateTimeLocal(c.PostForm("last_cleaned_date"))
	activeFrom := parseDateTimeLocal(c.PostForm("active_from"))
	activeTo := parseDateTimeLocal(c.PostForm("active_to"))

	lastServiceMileage := nullableInt(c.PostForm("last_service_mileage"))
	currentMileage := nullableInt(c.PostForm("current_mileage"))
	currentFuelLevel := nullableInt(c.PostForm("current_fuel_level"))

	_, err := h.svc.UpdateVehicle(
		id,
		vehicleMakeID, vehicleModelID, vehicleTypeID, fuelTypeID, vehicleColorID,
		description, plateNumber, iuNumber, chassisNumber, engineNumber, bootSpace,
		numSeats, carParkID, assetOwnerID, vehicleStatusID,
		lastServiceDate, lastCleanedDate, activeFrom, activeTo,
		lastServiceMileage, currentMileage, currentFuelLevel,
		userEmail,
	)
	if err != nil {
		slog.Error("failed to update vehicle", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to update vehicle")
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicles")
}

func (h *VehicleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteVehicle(id); err != nil {
		slog.Error("failed to delete vehicle", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete vehicle")
		return
	}

	c.Status(http.StatusOK)
}

func (h *VehicleHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tz := c.DefaultQuery("tz", "UTC")

	v, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch vehicle", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch vehicle")
		return
	}
	if v == nil {
		c.String(http.StatusNotFound, "Vehicle not found")
		return
	}

	c.HTML(http.StatusOK, "pages/vehicles/view.html", gin.H{
		"title":   "View Vehicle",
		"vehicle": v,
		"tz":      tz,
	})
}

func (h *VehicleHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	v, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch vehicle", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch vehicle")
		return
	}
	if v == nil {
		c.String(http.StatusNotFound, "Vehicle not found")
		return
	}

	name := fmt.Sprintf("Vehicle #%d", v.ID)
	if v.PlateNumber != nil && *v.PlateNumber != "" {
		name = *v.PlateNumber
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      name,
		"DeleteURL": "/vehicles/" + idStr,
		"RowID":     "vehicle-row-" + idStr,
	})
}

func nullableInt(s string) *int {
	if s == "" {
		return nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &n
}
