package web

import (
	"log/slog"
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type CondoPostalCodeHandler struct {
	svc      vehicle.CondoPostalCodeService
	condoSvc vehicle.CondoService
}

func NewCondoPostalCodeHandler(svc vehicle.CondoPostalCodeService, condoSvc vehicle.CondoService) *CondoPostalCodeHandler {
	return &CondoPostalCodeHandler{svc: svc, condoSvc: condoSvc}
}

func (h *CondoPostalCodeHandler) List(c *gin.Context) {
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
		slog.Error("failed to list condo postal codes", "error", err)
		c.String(http.StatusInternalServerError, "Failed to list condo postal codes")
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

	c.HTML(http.StatusOK, "pages/condo_postal_codes/index.html", gin.H{
		"title":            "Condominium Postal Codes",
		"condoPostalCodes": records,
		"timezone":         tz,
		"timezones":        Timezones,
		"currentPage":      page,
		"pageSize":         pageSize,
		"pageSizeOptions":  []int{5, 10, 20, 50},
		"totalPages":       totalPages,
		"total":            total,
		"start":            start,
		"end":              end,
		"sortBy":           sortBy,
		"sortOrder":        sortOrder,
		"pageWindow":       paginationWindow(page, totalPages),
	})
}

func (h *CondoPostalCodeHandler) CreateForm(c *gin.Context) {
	condos, err := h.condoSvc.ListAll()
	if err != nil {
		slog.Error("failed to list condos for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	c.HTML(http.StatusOK, "pages/condo_postal_codes/form.html", gin.H{
		"title":  "Add Condominium Postal Code",
		"action": "/condo-postal-codes",
		"condos": condos,
	})
}

func (h *CondoPostalCodeHandler) Create(c *gin.Context) {
	condoID, _ := strconv.ParseInt(c.PostForm("condo_id"), 10, 64)
	postalCode := c.PostForm("postal_code")

	_, err := h.svc.CreateCondoPostalCode(condoID, postalCode)
	if err != nil {
		slog.Error("failed to create condo postal code", "error", err)
		c.String(http.StatusInternalServerError, "Failed to create condo postal code")
		return
	}

	c.Redirect(http.StatusSeeOther, "/condo-postal-codes")
}

func (h *CondoPostalCodeHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	record, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch condo postal code", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch condo postal code")
		return
	}
	if record == nil {
		c.String(http.StatusNotFound, "Condo postal code not found")
		return
	}

	condos, err := h.condoSvc.ListAll()
	if err != nil {
		slog.Error("failed to list condos for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	c.HTML(http.StatusOK, "pages/condo_postal_codes/form.html", gin.H{
		"title":           "Edit Condominium Postal Code",
		"action":          "/condo-postal-codes/" + idStr,
		"condoPostalCode": record,
		"condos":          condos,
	})
}

func (h *CondoPostalCodeHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	condoID, _ := strconv.ParseInt(c.PostForm("condo_id"), 10, 64)
	postalCode := c.PostForm("postal_code")

	_, err := h.svc.UpdateCondoPostalCode(id, condoID, postalCode)
	if err != nil {
		slog.Error("failed to update condo postal code", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to update condo postal code")
		return
	}

	c.Redirect(http.StatusSeeOther, "/condo-postal-codes")
}

func (h *CondoPostalCodeHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteCondoPostalCode(id); err != nil {
		slog.Error("failed to delete condo postal code", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete condo postal code")
		return
	}

	c.Status(http.StatusOK)
}

func (h *CondoPostalCodeHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	record, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch condo postal code", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch condo postal code")
		return
	}
	if record == nil {
		c.String(http.StatusNotFound, "Condo postal code not found")
		return
	}

	c.HTML(http.StatusOK, "pages/condo_postal_codes/view.html", gin.H{
		"title":           "View Condominium Postal Code",
		"condoPostalCode": record,
	})
}

func (h *CondoPostalCodeHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	record, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch condo postal code", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch condo postal code")
		return
	}
	if record == nil {
		c.String(http.StatusNotFound, "Condo postal code not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      record.PostalCode,
		"DeleteURL": "/condo-postal-codes/" + idStr,
		"RowID":     "condo-postal-code-row-" + idStr,
	})
}
