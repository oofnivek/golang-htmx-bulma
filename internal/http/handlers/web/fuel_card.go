package web

import (
	"log/slog"
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type FuelCardHandler struct {
	svc   vehicle.FuelCardService
	fcSvc vehicle.FuelCompanyService
	vSvc  vehicle.VehicleService
}

func NewFuelCardHandler(svc vehicle.FuelCardService, fcSvc vehicle.FuelCompanyService, vSvc vehicle.VehicleService) *FuelCardHandler {
	return &FuelCardHandler{svc: svc, fcSvc: fcSvc, vSvc: vSvc}
}

func (h *FuelCardHandler) List(c *gin.Context) {
	tz := c.Query("tz")
	if tz == "" {
		tz = "UTC"
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	cards, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		slog.Error("failed to list fuel cards", "error", err)
		c.String(http.StatusInternalServerError, "Failed to list fuel cards")
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

	c.HTML(http.StatusOK, "pages/fuel_cards/index.html", gin.H{
		"title":           "Fuel Cards",
		"cards":           cards,
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

func (h *FuelCardHandler) CreateForm(c *gin.Context) {
	companies, err := h.fcSvc.ListAll()
	if err != nil {
		slog.Error("failed to list fuel companies for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	vehicles, err := h.vSvc.ListAll()
	if err != nil {
		slog.Error("failed to list vehicles for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	c.HTML(http.StatusOK, "pages/fuel_cards/form.html", gin.H{
		"title":     "Add Fuel Card",
		"action":    "/fuel-cards",
		"companies": companies,
		"vehicles":  vehicles,
	})
}

func (h *FuelCardHandler) Create(c *gin.Context) {
	cardNo := c.PostForm("card_no")
	fuelCompanyID, _ := strconv.ParseInt(c.PostForm("fuel_company_id"), 10, 64)
	pinNumber := c.PostForm("pin_number")
	status := c.PostForm("status") == "on"
	userEmail := c.GetString("user_email")

	var vehicleID *int64
	if v := c.PostForm("vehicle_id"); v != "" && v != "0" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			vehicleID = &id
		}
	}

	_, err := h.svc.CreateFuelCard(cardNo, fuelCompanyID, pinNumber, vehicleID, status, userEmail)
	if err != nil {
		slog.Error("failed to create fuel card", "error", err)
		c.String(http.StatusInternalServerError, "Failed to create fuel card")
		return
	}

	c.Redirect(http.StatusSeeOther, "/fuel-cards")
}

func (h *FuelCardHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	card, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch fuel card", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch fuel card")
		return
	}
	if card == nil {
		c.String(http.StatusNotFound, "Fuel card not found")
		return
	}

	companies, err := h.fcSvc.ListAll()
	if err != nil {
		slog.Error("failed to list fuel companies for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	vehicles, err := h.vSvc.ListAll()
	if err != nil {
		slog.Error("failed to list vehicles for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	c.HTML(http.StatusOK, "pages/fuel_cards/form.html", gin.H{
		"title":     "Edit Fuel Card",
		"action":    "/fuel-cards/" + idStr,
		"card":      card,
		"companies": companies,
		"vehicles":  vehicles,
	})
}

func (h *FuelCardHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	cardNo := c.PostForm("card_no")
	fuelCompanyID, _ := strconv.ParseInt(c.PostForm("fuel_company_id"), 10, 64)
	pinNumber := c.PostForm("pin_number")
	status := c.PostForm("status") == "on"
	userEmail := c.GetString("user_email")

	var vehicleID *int64
	if v := c.PostForm("vehicle_id"); v != "" && v != "0" {
		vid, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			vehicleID = &vid
		}
	}

	_, err := h.svc.UpdateFuelCard(id, cardNo, fuelCompanyID, pinNumber, vehicleID, status, userEmail)
	if err != nil {
		slog.Error("failed to update fuel card", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to update fuel card")
		return
	}

	c.Redirect(http.StatusSeeOther, "/fuel-cards")
}

func (h *FuelCardHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteFuelCard(id); err != nil {
		slog.Error("failed to delete fuel card", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete fuel card")
		return
	}

	c.Status(http.StatusOK)
}

func (h *FuelCardHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tz := c.DefaultQuery("tz", "UTC")

	card, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch fuel card", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch fuel card")
		return
	}
	if card == nil {
		c.String(http.StatusNotFound, "Fuel card not found")
		return
	}

	c.HTML(http.StatusOK, "pages/fuel_cards/view.html", gin.H{
		"title": "View Fuel Card",
		"card":  card,
		"tz":    tz,
	})
}

func (h *FuelCardHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	card, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch fuel card", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch fuel card")
		return
	}
	if card == nil {
		c.String(http.StatusNotFound, "Fuel card not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      card.CardNo,
		"DeleteURL": "/fuel-cards/" + idStr,
		"RowID":     "fuel-card-row-" + idStr,
	})
}
