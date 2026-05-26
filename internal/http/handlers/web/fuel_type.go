package web

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type FuelTypeHandler struct {
	svc vehicle.FuelTypeService
}

func NewFuelTypeHandler(svc vehicle.FuelTypeService) *FuelTypeHandler {
	return &FuelTypeHandler{svc: svc}
}

func (h *FuelTypeHandler) List(c *gin.Context) {
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

	c.HTML(http.StatusOK, "pages/fuel_types/index.html", gin.H{
		"title":           "Fuel Types",
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

func (h *FuelTypeHandler) CreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/fuel_types/form.html", gin.H{
		"title":  "Add Fuel Type",
		"action": "/fuel-types",
	})
}

func (h *FuelTypeHandler) Create(c *gin.Context) {
	name := c.PostForm("name")
	status := c.PostForm("status") == "on"
	userEmail := c.GetString("user_email")

	_, err := h.svc.CreateFuelType(name, status, userEmail)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/fuel-types")
}

func (h *FuelTypeHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	f, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if f == nil {
		c.String(http.StatusNotFound, "Fuel type not found")
		return
	}

	c.HTML(http.StatusOK, "pages/fuel_types/form.html", gin.H{
		"title":  "Edit Fuel Type",
		"action": "/fuel-types/" + idStr,
		"fuel":   f,
	})
}

func (h *FuelTypeHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	name := c.PostForm("name")
	status := c.PostForm("status") == "on"
	userEmail := c.GetString("user_email")

	_, err := h.svc.UpdateFuelType(id, name, status, userEmail)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/fuel-types")
}

func (h *FuelTypeHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteFuelType(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *FuelTypeHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tz := c.DefaultQuery("tz", "UTC")

	f, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if f == nil {
		c.String(http.StatusNotFound, "Fuel type not found")
		return
	}

	c.HTML(http.StatusOK, "pages/fuel_types/view.html", gin.H{
		"title": "View Fuel Type",
		"fuel":  f,
		"tz":    tz,
	})
}

func (h *FuelTypeHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	f, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if f == nil {
		c.String(http.StatusNotFound, "Fuel type not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      f.Name,
		"DeleteURL": "/fuel-types/" + idStr,
		"RowID":     "fuel-row-" + idStr,
	})
}
