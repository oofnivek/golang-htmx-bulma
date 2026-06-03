package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type CondoAcknowledgementHandler struct {
	svc vehicle.CondoAcknowledgementService
}

func NewCondoAcknowledgementHandler(svc vehicle.CondoAcknowledgementService) *CondoAcknowledgementHandler {
	return &CondoAcknowledgementHandler{svc: svc}
}

func (h *CondoAcknowledgementHandler) List(c *gin.Context) {
	tz := c.Query("tz")
	if tz == "" {
		tz = "UTC"
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	records, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		slog.Error("failed to list condo acknowledgements", "error", err)
		c.String(http.StatusInternalServerError, "Failed to list condo acknowledgements")
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

	c.HTML(http.StatusOK, "pages/condo_acknowledgements/index.html", gin.H{
		"title":                 "Condominium Acknowledgements",
		"condoAcknowledgements": records,
		"timezone":              tz,
		"timezones":             Timezones,
		"currentPage":           page,
		"pageSize":              pageSize,
		"pageSizeOptions":       []int{5, 10, 20, 50},
		"totalPages":            totalPages,
		"total":                 total,
		"start":                 start,
		"end":                   end,
		"sortBy":                sortBy,
		"sortOrder":             sortOrder,
		"pageWindow":            paginationWindow(page, totalPages),
	})
}

func (h *CondoAcknowledgementHandler) CreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/condo_acknowledgements/form.html", gin.H{
		"title":  "Add Condominium Acknowledgement",
		"action": "/condo-acknowledgements",
	})
}

func (h *CondoAcknowledgementHandler) Create(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.PostForm("user_id"), 10, 64)

	_, err := h.svc.CreateCondoAcknowledgement(userID)
	if err != nil {
		slog.Error("failed to create condo acknowledgement", "user_id", userID, "error", err)
		c.String(http.StatusInternalServerError, "Failed to create condo acknowledgement")
		return
	}

	c.Redirect(http.StatusSeeOther, "/condo-acknowledgements")
}

func (h *CondoAcknowledgementHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tz := c.DefaultQuery("tz", "UTC")

	record, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch condo acknowledgement", "user_id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch condo acknowledgement")
		return
	}
	if record == nil {
		c.String(http.StatusNotFound, "Condo acknowledgement not found")
		return
	}

	c.HTML(http.StatusOK, "pages/condo_acknowledgements/view.html", gin.H{
		"title":                "View Condominium Acknowledgement",
		"condoAcknowledgement": record,
		"tz":                   tz,
	})
}

func (h *CondoAcknowledgementHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteCondoAcknowledgement(id); err != nil {
		slog.Error("failed to delete condo acknowledgement", "user_id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete condo acknowledgement")
		return
	}

	c.Status(http.StatusOK)
}

func (h *CondoAcknowledgementHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	record, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch condo acknowledgement", "user_id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch condo acknowledgement")
		return
	}
	if record == nil {
		c.String(http.StatusNotFound, "Condo acknowledgement not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      fmt.Sprintf("User #%d", record.UserID),
		"DeleteURL": "/condo-acknowledgements/" + idStr,
		"RowID":     "condo-acknowledgement-row-" + idStr,
	})
}
