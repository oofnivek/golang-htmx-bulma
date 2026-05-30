package web

import (
	"log/slog"
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type CarParkLotHandler struct {
	svc   vehicle.CarParkLotService
	cpSvc vehicle.CarParkService
}

func NewCarParkLotHandler(svc vehicle.CarParkLotService, cpSvc vehicle.CarParkService) *CarParkLotHandler {
	return &CarParkLotHandler{svc: svc, cpSvc: cpSvc}
}

func (h *CarParkLotHandler) List(c *gin.Context) {
	tz := c.Query("tz")
	if tz == "" {
		tz = "UTC"
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	lots, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		slog.Error("failed to list car park lots", "error", err)
		c.String(http.StatusInternalServerError, "Failed to list car park lots")
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

	c.HTML(http.StatusOK, "pages/car_park_lots/index.html", gin.H{
		"title":           "Car Park Lots",
		"lots":            lots,
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

func (h *CarParkLotHandler) CreateForm(c *gin.Context) {
	parks, err := h.cpSvc.ListAll()
	if err != nil {
		slog.Error("failed to list car parks for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	c.HTML(http.StatusOK, "pages/car_park_lots/form.html", gin.H{
		"title":  "Add Car Park Lot",
		"action": "/car-park-lots",
		"parks":  parks,
	})
}

func (h *CarParkLotHandler) Create(c *gin.Context) {
	carParkID, _ := strconv.ParseInt(c.PostForm("car_park_id"), 10, 64)
	lotNumber := c.PostForm("lot_number")
	level := c.PostForm("level")
	status := c.PostForm("status") == "on"
	userEmail := c.GetString("user_email")

	_, err := h.svc.CreateCarParkLot(carParkID, lotNumber, level, status, userEmail)
	if err != nil {
		slog.Error("failed to create car park lot", "error", err)
		c.String(http.StatusInternalServerError, "Failed to create car park lot")
		return
	}

	c.Redirect(http.StatusSeeOther, "/car-park-lots")
}

func (h *CarParkLotHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	l, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch car park lot", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch car park lot")
		return
	}
	if l == nil {
		c.String(http.StatusNotFound, "Car park lot not found")
		return
	}

	parks, err := h.cpSvc.ListAll()
	if err != nil {
		slog.Error("failed to list car parks for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	c.HTML(http.StatusOK, "pages/car_park_lots/form.html", gin.H{
		"title":  "Edit Car Park Lot",
		"action": "/car-park-lots/" + idStr,
		"lot":    l,
		"parks":  parks,
	})
}

func (h *CarParkLotHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	carParkID, _ := strconv.ParseInt(c.PostForm("car_park_id"), 10, 64)
	lotNumber := c.PostForm("lot_number")
	level := c.PostForm("level")
	status := c.PostForm("status") == "on"
	userEmail := c.GetString("user_email")

	_, err := h.svc.UpdateCarParkLot(id, carParkID, lotNumber, level, status, userEmail)
	if err != nil {
		slog.Error("failed to update car park lot", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to update car park lot")
		return
	}

	c.Redirect(http.StatusSeeOther, "/car-park-lots")
}

func (h *CarParkLotHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteCarParkLot(id); err != nil {
		slog.Error("failed to delete car park lot", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete car park lot")
		return
	}

	c.Status(http.StatusOK)
}

func (h *CarParkLotHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tz := c.DefaultQuery("tz", "UTC")

	l, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch car park lot", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch car park lot")
		return
	}
	if l == nil {
		c.String(http.StatusNotFound, "Car park lot not found")
		return
	}

	c.HTML(http.StatusOK, "pages/car_park_lots/view.html", gin.H{
		"title": "View Car Park Lot",
		"lot":   l,
		"tz":    tz,
	})
}

func (h *CarParkLotHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	l, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch car park lot", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch car park lot")
		return
	}
	if l == nil {
		c.String(http.StatusNotFound, "Car park lot not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      l.LotNumber,
		"DeleteURL": "/car-park-lots/" + idStr,
		"RowID":     "car-park-lot-row-" + idStr,
	})
}
