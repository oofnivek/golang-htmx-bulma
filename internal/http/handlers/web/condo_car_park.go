package web

import (
	"log/slog"
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type CondoCarParkHandler struct {
	svc      vehicle.CondoCarParkService
	condoSvc vehicle.CondoService
	cpSvc    vehicle.CarParkService
}

func NewCondoCarParkHandler(svc vehicle.CondoCarParkService, condoSvc vehicle.CondoService, cpSvc vehicle.CarParkService) *CondoCarParkHandler {
	return &CondoCarParkHandler{svc: svc, condoSvc: condoSvc, cpSvc: cpSvc}
}

func (h *CondoCarParkHandler) List(c *gin.Context) {
	tz := c.Query("tz")
	if tz == "" {
		tz = "UTC"
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	records, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		slog.Error("failed to list condo car parks", "error", err)
		c.String(http.StatusInternalServerError, "Failed to list condo car parks")
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

	c.HTML(http.StatusOK, "pages/condo_car_parks/index.html", gin.H{
		"title":           "Condominium Car Parks",
		"condoCarParks":   records,
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

func (h *CondoCarParkHandler) CreateForm(c *gin.Context) {
	condos, err := h.condoSvc.ListAll()
	if err != nil {
		slog.Error("failed to list condos for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	parks, err := h.cpSvc.ListAll()
	if err != nil {
		slog.Error("failed to list car parks for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	c.HTML(http.StatusOK, "pages/condo_car_parks/form.html", gin.H{
		"title":  "Add Condominium Car Park",
		"action": "/condo-car-parks",
		"condos": condos,
		"parks":  parks,
	})
}

func (h *CondoCarParkHandler) Create(c *gin.Context) {
	condoID, _ := strconv.ParseInt(c.PostForm("condo_id"), 10, 64)
	carParkID, _ := strconv.ParseInt(c.PostForm("car_park_id"), 10, 64)

	_, err := h.svc.CreateCondoCarPark(condoID, carParkID)
	if err != nil {
		slog.Error("failed to create condo car park", "error", err)
		c.String(http.StatusInternalServerError, "Failed to create condo car park")
		return
	}

	c.Redirect(http.StatusSeeOther, "/condo-car-parks")
}

func (h *CondoCarParkHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	record, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch condo car park", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch condo car park")
		return
	}
	if record == nil {
		c.String(http.StatusNotFound, "Condo car park not found")
		return
	}

	condos, err := h.condoSvc.ListAll()
	if err != nil {
		slog.Error("failed to list condos for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	parks, err := h.cpSvc.ListAll()
	if err != nil {
		slog.Error("failed to list car parks for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	c.HTML(http.StatusOK, "pages/condo_car_parks/form.html", gin.H{
		"title":        "Edit Condominium Car Park",
		"action":       "/condo-car-parks/" + idStr,
		"condoCarPark": record,
		"condos":       condos,
		"parks":        parks,
	})
}

func (h *CondoCarParkHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	condoID, _ := strconv.ParseInt(c.PostForm("condo_id"), 10, 64)
	carParkID, _ := strconv.ParseInt(c.PostForm("car_park_id"), 10, 64)

	_, err := h.svc.UpdateCondoCarPark(id, condoID, carParkID)
	if err != nil {
		slog.Error("failed to update condo car park", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to update condo car park")
		return
	}

	c.Redirect(http.StatusSeeOther, "/condo-car-parks")
}

func (h *CondoCarParkHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteCondoCarPark(id); err != nil {
		slog.Error("failed to delete condo car park", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete condo car park")
		return
	}

	c.Status(http.StatusOK)
}

func (h *CondoCarParkHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	record, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch condo car park", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch condo car park")
		return
	}
	if record == nil {
		c.String(http.StatusNotFound, "Condo car park not found")
		return
	}

	c.HTML(http.StatusOK, "pages/condo_car_parks/view.html", gin.H{
		"title":        "View Condominium Car Park",
		"condoCarPark": record,
	})
}

func (h *CondoCarParkHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	record, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch condo car park", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch condo car park")
		return
	}
	if record == nil {
		c.String(http.StatusNotFound, "Condo car park not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      record.CondoName + " / " + record.CarParkName,
		"DeleteURL": "/condo-car-parks/" + idStr,
		"RowID":     "condo-car-park-row-" + idStr,
	})
}
