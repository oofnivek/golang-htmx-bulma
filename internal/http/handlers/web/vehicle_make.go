package web

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type VehicleMakeHandler struct {
	svc vehicle.VehicleMakeService
}

func NewVehicleMakeHandler(svc vehicle.VehicleMakeService) *VehicleMakeHandler {
	return &VehicleMakeHandler{svc: svc}
}

func (h *VehicleMakeHandler) List(c *gin.Context) {
	tz := c.Query("tz")
	if tz == "" {
		tz = "UTC"
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	makes, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	timezones := Timezones

	totalPages := (total + pageSize - 1) / pageSize

	c.HTML(http.StatusOK, "pages/vehicle_makes/index.html", gin.H{
		"title":           "Vehicle Makes",
		"makes":           makes,
		"timezone":        tz,
		"timezones":       timezones,
		"currentPage":     page,
		"pageSize":        pageSize,
		"pageSizeOptions": []int{5, 10, 20, 50},
		"totalPages":      totalPages,
		"total":           total,
		"sortBy":          sortBy,
		"sortOrder":       sortOrder,
	})
}

func (h *VehicleMakeHandler) CreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/vehicle_makes/form.html", gin.H{
		"title":  "Add Vehicle Make",
		"action": "/vehicle-makes",
	})
}

func (h *VehicleMakeHandler) Create(c *gin.Context) {
	name := c.PostForm("name")
	statusStr := c.PostForm("status")
	status := statusStr == "on"

	userEmail := c.GetString("user_email")
	_, err := h.svc.CreateMake(name, status, userEmail) // Use authenticated user email
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicle-makes")
}

func (h *VehicleMakeHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	make, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if make == nil {
		c.String(http.StatusNotFound, "Make not found")
		return
	}

	c.HTML(http.StatusOK, "pages/vehicle_makes/form.html", gin.H{
		"title":  "Edit Vehicle Make",
		"action": "/vehicle-makes/" + idStr,
		"make":   make,
	})
}

func (h *VehicleMakeHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	name := c.PostForm("name")
	statusStr := c.PostForm("status")
	status := statusStr == "on"

	userEmail := c.GetString("user_email")
	_, err := h.svc.UpdateMake(id, name, status, userEmail)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicle-makes")
}

func (h *VehicleMakeHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	err := h.svc.DeleteMake(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *VehicleMakeHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tz := c.DefaultQuery("tz", "UTC")

	m, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if m == nil {
		c.String(http.StatusNotFound, "Make not found")
		return
	}

	c.HTML(http.StatusOK, "pages/vehicle_makes/view.html", gin.H{
		"title": "View Vehicle Make",
		"make":  m,
		"tz":    tz,
	})
}

func (h *VehicleMakeHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	make, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if make == nil {
		c.String(http.StatusNotFound, "Make not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      make.Name,
		"DeleteURL": "/vehicle-makes/" + idStr,
		"RowID":     "make-row-" + idStr,
	})
}
