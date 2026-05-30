package web

import (
	"log/slog"
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type VehicleModelHandler struct {
	svc     vehicle.VehicleModelService
	typeSvc vehicle.VehicleTypeService
	makeSvc vehicle.VehicleMakeService
}

func NewVehicleModelHandler(svc vehicle.VehicleModelService, typeSvc vehicle.VehicleTypeService, makeSvc vehicle.VehicleMakeService) *VehicleModelHandler {
	return &VehicleModelHandler{svc: svc, typeSvc: typeSvc, makeSvc: makeSvc}
}

func (h *VehicleModelHandler) List(c *gin.Context) {
	tz := c.Query("tz")
	if tz == "" {
		tz = "UTC"
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	models, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	c.HTML(http.StatusOK, "pages/vehicle_models/index.html", gin.H{
		"title":           "Vehicle Models",
		"models":          models,
		"timezone":        tz,
		"timezones":       Timezones,
		"currentPage":     page,
		"pageSize":        pageSize,
		"pageSizeOptions": []int{5, 10, 20, 50},
		"totalPages":      totalPages,
"pageWindow":      paginationWindow(page, totalPages),
		"total":           total,
		"sortBy":          sortBy,
		"sortOrder":       sortOrder,
	})
}

func (h *VehicleModelHandler) CreateForm(c *gin.Context) {
	types, err := h.typeSvc.ListAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	makes, err := h.makeSvc.ListAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "pages/vehicle_models/form.html", gin.H{
		"title":        "Add Vehicle Model",
		"action":       "/vehicle-models",
		"vehicleTypes": types,
		"vehicleMakes": makes,
	})
}

func (h *VehicleModelHandler) Create(c *gin.Context) {
	vehicleTypeID, _ := strconv.ParseInt(c.PostForm("vehicle_type_id"), 10, 64)
	vehicleMakeID, _ := strconv.ParseInt(c.PostForm("vehicle_make_id"), 10, 64)
	name := c.PostForm("name")
	status := c.PostForm("status") == "on"

	userEmail := c.GetString("user_email")
	_, err := h.svc.CreateModel(vehicleTypeID, vehicleMakeID, name, status, userEmail)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicle-models")
}

func (h *VehicleModelHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	m, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if m == nil {
		c.String(http.StatusNotFound, "Model not found")
		return
	}

	types, err := h.typeSvc.ListAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	makes, err := h.makeSvc.ListAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "pages/vehicle_models/form.html", gin.H{
		"title":          "Edit Vehicle Model",
		"action":         "/vehicle-models/" + idStr,
		"model":          m,
		"vehicleTypes":   types,
		"vehicleMakes":   makes,
		"currentTypeID":  m.VehicleTypeID,
		"currentMakeID":  m.VehicleMakeID,
	})
}

func (h *VehicleModelHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	vehicleTypeID, _ := strconv.ParseInt(c.PostForm("vehicle_type_id"), 10, 64)
	vehicleMakeID, _ := strconv.ParseInt(c.PostForm("vehicle_make_id"), 10, 64)
	name := c.PostForm("name")
	status := c.PostForm("status") == "on"

	userEmail := c.GetString("user_email")
	_, err := h.svc.UpdateModel(id, vehicleTypeID, vehicleMakeID, name, status, userEmail)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicle-models")
}

func (h *VehicleModelHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteModel(id); err != nil {
		slog.Error("failed to delete vehicle model", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete vehicle model")
		return
	}

	c.Status(http.StatusOK)
}

func (h *VehicleModelHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tz := c.DefaultQuery("tz", "UTC")

	m, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if m == nil {
		c.String(http.StatusNotFound, "Vehicle model not found")
		return
	}

	c.HTML(http.StatusOK, "pages/vehicle_models/view.html", gin.H{
		"title": "View Vehicle Model",
		"model": m,
		"tz":    tz,
	})
}

func (h *VehicleModelHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	m, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if m == nil {
		c.String(http.StatusNotFound, "Model not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      m.Name,
		"DeleteURL": "/vehicle-models/" + idStr,
		"RowID":     "model-row-" + idStr,
	})
}
