package web

import (
	"log/slog"
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/pkg/status"
	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type CondoHandler struct{ svc vehicle.CondoService }

func NewCondoHandler(svc vehicle.CondoService) *CondoHandler {
	return &CondoHandler{svc: svc}
}

func (h *CondoHandler) List(c *gin.Context) {
	tz := c.Query("tz")
	if tz == "" {
		tz = "UTC"
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	condos, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		slog.Error("failed to list condos", "error", err)
		c.String(http.StatusInternalServerError, "Failed to list condos")
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

	c.HTML(http.StatusOK, "pages/condos/index.html", gin.H{
		"title":           "Condominiums",
		"condos":          condos,
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

func (h *CondoHandler) CreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/condos/form.html", gin.H{
		"title":  "Add Condominium",
		"action": "/condos",
	})
}

func (h *CondoHandler) Create(c *gin.Context) {
	name := c.PostForm("name")
	stInt, _ := strconv.Atoi(c.PostForm("status"))
	st := status.Status(stInt)
	mcstNumber := c.PostForm("mcst_number")
	mcstEmail := c.PostForm("mcst_email")
	address := c.PostForm("address")
	userEmail := c.GetString("user_email")

	_, err := h.svc.CreateCondo(name, st, mcstNumber, mcstEmail, address, userEmail)
	if err != nil {
		slog.Error("failed to create condo", "error", err)
		c.String(http.StatusInternalServerError, "Failed to create condo")
		return
	}

	c.Redirect(http.StatusSeeOther, "/condos")
}

func (h *CondoHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	condo, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch condo", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch condo")
		return
	}
	if condo == nil {
		c.String(http.StatusNotFound, "Condo not found")
		return
	}

	c.HTML(http.StatusOK, "pages/condos/form.html", gin.H{
		"title":  "Edit Condominium",
		"action": "/condos/" + idStr,
		"condo":  condo,
	})
}

func (h *CondoHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	name := c.PostForm("name")
	stInt, _ := strconv.Atoi(c.PostForm("status"))
	st := status.Status(stInt)
	mcstNumber := c.PostForm("mcst_number")
	mcstEmail := c.PostForm("mcst_email")
	address := c.PostForm("address")
	userEmail := c.GetString("user_email")

	_, err := h.svc.UpdateCondo(id, name, st, mcstNumber, mcstEmail, address, userEmail)
	if err != nil {
		slog.Error("failed to update condo", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to update condo")
		return
	}

	c.Redirect(http.StatusSeeOther, "/condos")
}

func (h *CondoHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteCondo(id); err != nil {
		slog.Error("failed to delete condo", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete condo")
		return
	}

	c.Status(http.StatusOK)
}

func (h *CondoHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tz := c.DefaultQuery("tz", "UTC")

	condo, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch condo", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch condo")
		return
	}
	if condo == nil {
		c.String(http.StatusNotFound, "Condo not found")
		return
	}

	c.HTML(http.StatusOK, "pages/condos/view.html", gin.H{
		"title": "View Condominium",
		"condo": condo,
		"tz":    tz,
	})
}

func (h *CondoHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	condo, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch condo", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch condo")
		return
	}
	if condo == nil {
		c.String(http.StatusNotFound, "Condo not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      condo.Name,
		"DeleteURL": "/condos/" + idStr,
		"RowID":     "condo-row-" + idStr,
	})
}
