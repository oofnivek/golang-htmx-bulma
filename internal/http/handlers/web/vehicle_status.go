package web

import (
	"log/slog"
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type VehicleStatusHandler struct {
	svc vehicle.VehicleStatusService
}

func NewVehicleStatusHandler(svc vehicle.VehicleStatusService) *VehicleStatusHandler {
	return &VehicleStatusHandler{svc: svc}
}

func (h *VehicleStatusHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	statuses, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	c.HTML(http.StatusOK, "pages/vehicle_statuses/index.html", gin.H{
		"title":           "Vehicle Statuses",
		"statuses":        statuses,
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

func (h *VehicleStatusHandler) CreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/vehicle_statuses/form.html", gin.H{
		"title":  "Add Vehicle Status",
		"action": "/vehicle-statuses",
	})
}

func (h *VehicleStatusHandler) Create(c *gin.Context) {
	substatus := c.PostForm("substatus")
	isActive := c.PostForm("is_active") == "on"

	_, err := h.svc.CreateStatus(substatus, isActive)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicle-statuses")
}

func (h *VehicleStatusHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	vs, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if vs == nil {
		c.String(http.StatusNotFound, "Status not found")
		return
	}

	c.HTML(http.StatusOK, "pages/vehicle_statuses/form.html", gin.H{
		"title":  "Edit Vehicle Status",
		"action": "/vehicle-statuses/" + idStr,
		"status": vs,
	})
}

func (h *VehicleStatusHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	substatus := c.PostForm("substatus")
	isActive := c.PostForm("is_active") == "on"

	_, err := h.svc.UpdateStatus(id, substatus, isActive)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicle-statuses")
}

func (h *VehicleStatusHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteStatus(id); err != nil {
		slog.Error("failed to delete vehicle status", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete vehicle status")
		return
	}

	c.Status(http.StatusOK)
}

func (h *VehicleStatusHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	vs, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if vs == nil {
		c.String(http.StatusNotFound, "Status not found")
		return
	}

	c.HTML(http.StatusOK, "pages/vehicle_statuses/view.html", gin.H{
		"title":  "View Vehicle Status",
		"status": vs,
	})
}

func (h *VehicleStatusHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	vs, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if vs == nil {
		c.String(http.StatusNotFound, "Status not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      vs.Substatus,
		"DeleteURL": "/vehicle-statuses/" + idStr,
		"RowID":     "status-row-" + idStr,
	})
}
