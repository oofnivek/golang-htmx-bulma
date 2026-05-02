package web

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/service"

	"github.com/gin-gonic/gin"
)

type VehicleColorHandler struct {
	svc service.VehicleColorService
}

func NewVehicleColorHandler(svc service.VehicleColorService) *VehicleColorHandler {
	return &VehicleColorHandler{svc: svc}
}

func (h *VehicleColorHandler) List(c *gin.Context) {
	tz := c.Query("tz")
	if tz == "" {
		tz = "UTC"
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")
	search := c.Query("search")

	colors, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder, search)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Common timezones
	timezones := []string{"UTC", "America/New_York", "America/Los_Angeles", "Europe/London", "Europe/Paris", "Asia/Tokyo", "Asia/Shanghai", "Asia/Singapore", "Australia/Sydney"}

	totalPages := (total + pageSize - 1) / pageSize

	c.HTML(http.StatusOK, "pages/vehicle_colors/index.html", gin.H{
		"title":           "Vehicle Colors",
		"colors":          colors,
		"timezone":        tz,
		"timezones":       timezones,
		"currentPage":     page,
		"pageSize":        pageSize,
		"pageSizeOptions": []int{5, 10, 20, 50},
		"totalPages":      totalPages,
		"total":           total,
		"sortBy":          sortBy,
		"sortOrder":       sortOrder,
		"search":          search,
	})
}

func (h *VehicleColorHandler) CreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/vehicle_colors/form.html", gin.H{
		"title":  "Add Vehicle Color",
		"action": "/vehicle-colors",
	})
}

func (h *VehicleColorHandler) Create(c *gin.Context) {
	name := c.PostForm("name")
	statusStr := c.PostForm("status")
	status := statusStr == "on" // Bulma checkbox or select? Let's assume checkbox for now.

	_, err := h.svc.CreateColor(name, status, "admin") // Hardcoded user for now
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicle-colors")
}

func (h *VehicleColorHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	color, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if color == nil {
		c.String(http.StatusNotFound, "Color not found")
		return
	}

	c.HTML(http.StatusOK, "pages/vehicle_colors/form.html", gin.H{
		"title":  "Edit Vehicle Color",
		"action": "/vehicle-colors/" + idStr,
		"color":  color,
	})
}

func (h *VehicleColorHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	name := c.PostForm("name")
	statusStr := c.PostForm("status")
	status := statusStr == "on"

	_, err := h.svc.UpdateColor(id, name, status, "admin")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicle-colors")
}

func (h *VehicleColorHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	err := h.svc.DeleteColor(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *VehicleColorHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	color, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if color == nil {
		c.String(http.StatusNotFound, "Color not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      color.Name,
		"DeleteURL": "/vehicle-colors/" + idStr,
		"RowID":     "color-row-" + idStr,
	})
}
