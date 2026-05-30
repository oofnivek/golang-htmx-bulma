package web

import (
	"log/slog"
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type EstateHandler struct {
	svc vehicle.EstateService
}

func NewEstateHandler(svc vehicle.EstateService) *EstateHandler {
	return &EstateHandler{svc: svc}
}

func (h *EstateHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	estates, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		slog.Error("failed to list estates", "error", err)
		c.String(http.StatusInternalServerError, "Failed to list estates")
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

	c.HTML(http.StatusOK, "pages/estates/index.html", gin.H{
		"title":           "Estates",
		"estates":         estates,
		"currentPage":     page,
		"pageSize":        pageSize,
		"pageSizeOptions": []int{5, 10, 20, 50},
		"totalPages":      totalPages,
		"total":           total,
		"start":           start,
		"end":             end,
		"sortBy":          sortBy,
		"sortOrder":       sortOrder,
	})
}

func (h *EstateHandler) CreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/estates/form.html", gin.H{
		"title":  "Add Estate",
		"action": "/estates",
	})
}

func (h *EstateHandler) Create(c *gin.Context) {
	name := c.PostForm("name")

	_, err := h.svc.CreateEstate(name)
	if err != nil {
		slog.Error("failed to create estate", "error", err)
		c.String(http.StatusInternalServerError, "Failed to create estate")
		return
	}

	c.Redirect(http.StatusSeeOther, "/estates")
}

func (h *EstateHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	e, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch estate", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch estate")
		return
	}
	if e == nil {
		c.String(http.StatusNotFound, "Estate not found")
		return
	}

	c.HTML(http.StatusOK, "pages/estates/form.html", gin.H{
		"title":  "Edit Estate",
		"action": "/estates/" + idStr,
		"estate": e,
	})
}

func (h *EstateHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	name := c.PostForm("name")

	_, err := h.svc.UpdateEstate(id, name)
	if err != nil {
		slog.Error("failed to update estate", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to update estate")
		return
	}

	c.Redirect(http.StatusSeeOther, "/estates")
}

func (h *EstateHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteEstate(id); err != nil {
		slog.Error("failed to delete estate", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete estate")
		return
	}

	c.Status(http.StatusOK)
}

func (h *EstateHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	e, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch estate", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch estate")
		return
	}
	if e == nil {
		c.String(http.StatusNotFound, "Estate not found")
		return
	}

	c.HTML(http.StatusOK, "pages/estates/view.html", gin.H{
		"title":  "View Estate",
		"estate": e,
	})
}

func (h *EstateHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	e, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch estate", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch estate")
		return
	}
	if e == nil {
		c.String(http.StatusNotFound, "Estate not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      e.Name,
		"DeleteURL": "/estates/" + idStr,
		"RowID":     "estate-row-" + idStr,
	})
}
