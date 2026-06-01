package web

import (
	"log/slog"
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type RegionalInfoHandler struct {
	svc      vehicle.RegionalInfoService
	estateSvc vehicle.EstateService
}

func NewRegionalInfoHandler(svc vehicle.RegionalInfoService, estateSvc vehicle.EstateService) *RegionalInfoHandler {
	return &RegionalInfoHandler{svc: svc, estateSvc: estateSvc}
}

func (h *RegionalInfoHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "postal_code")
	sortOrder := c.DefaultQuery("sortOrder", "asc")

	items, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		slog.Error("failed to list regional infos", "error", err)
		c.String(http.StatusInternalServerError, "Failed to list regional infos")
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

	c.HTML(http.StatusOK, "pages/regional_infos/index.html", gin.H{
		"title":           "Regional Info",
		"items":           items,
		"currentPage":     page,
		"pageSize":        pageSize,
		"pageSizeOptions": []int{5, 10, 20, 50},
		"totalPages":      totalPages,
		"pageWindow":      paginationWindow(page, totalPages),
		"total":           total,
		"start":           start,
		"end":             end,
		"sortBy":          sortBy,
		"sortOrder":       sortOrder,
	})
}

func (h *RegionalInfoHandler) CreateForm(c *gin.Context) {
	estates, err := h.estateSvc.ListAll()
	if err != nil {
		slog.Error("failed to list estates for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	c.HTML(http.StatusOK, "pages/regional_infos/form.html", gin.H{
		"title":   "Add Regional Info",
		"action":  "/regional-infos",
		"estates": estates,
	})
}

func (h *RegionalInfoHandler) Create(c *gin.Context) {
	postalCode := c.PostForm("postal_code")
	region := c.PostForm("region")
	estateID, _ := strconv.ParseInt(c.PostForm("estate_id"), 10, 64)

	_, err := h.svc.CreateRegionalInfo(postalCode, region, estateID)
	if err != nil {
		slog.Error("failed to create regional info", "postal_code", postalCode, "error", err)
		c.String(http.StatusInternalServerError, "Failed to create regional info")
		return
	}

	c.Redirect(http.StatusSeeOther, "/regional-infos")
}

func (h *RegionalInfoHandler) EditForm(c *gin.Context) {
	postalCode := c.Param("id")

	ri, err := h.svc.FindByID(postalCode)
	if err != nil {
		slog.Error("failed to fetch regional info", "postal_code", postalCode, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch regional info")
		return
	}
	if ri == nil {
		c.String(http.StatusNotFound, "Regional info not found")
		return
	}

	estates, err := h.estateSvc.ListAll()
	if err != nil {
		slog.Error("failed to list estates for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	c.HTML(http.StatusOK, "pages/regional_infos/form.html", gin.H{
		"title":   "Edit Regional Info",
		"action":  "/regional-infos/" + postalCode,
		"item":    ri,
		"estates": estates,
	})
}

func (h *RegionalInfoHandler) Update(c *gin.Context) {
	postalCode := c.Param("id")
	region := c.PostForm("region")
	estateID, _ := strconv.ParseInt(c.PostForm("estate_id"), 10, 64)

	_, err := h.svc.UpdateRegionalInfo(postalCode, region, estateID)
	if err != nil {
		slog.Error("failed to update regional info", "postal_code", postalCode, "error", err)
		c.String(http.StatusInternalServerError, "Failed to update regional info")
		return
	}

	c.Redirect(http.StatusSeeOther, "/regional-infos")
}

func (h *RegionalInfoHandler) Delete(c *gin.Context) {
	postalCode := c.Param("id")

	if err := h.svc.DeleteRegionalInfo(postalCode); err != nil {
		slog.Error("failed to delete regional info", "postal_code", postalCode, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete regional info")
		return
	}

	c.Status(http.StatusOK)
}

func (h *RegionalInfoHandler) View(c *gin.Context) {
	postalCode := c.Param("id")

	ri, err := h.svc.FindByID(postalCode)
	if err != nil {
		slog.Error("failed to fetch regional info", "postal_code", postalCode, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch regional info")
		return
	}
	if ri == nil {
		c.String(http.StatusNotFound, "Regional info not found")
		return
	}

	c.HTML(http.StatusOK, "pages/regional_infos/view.html", gin.H{
		"title": "View Regional Info",
		"item":  ri,
	})
}

func (h *RegionalInfoHandler) DeleteConfirm(c *gin.Context) {
	postalCode := c.Param("id")

	ri, err := h.svc.FindByID(postalCode)
	if err != nil {
		slog.Error("failed to fetch regional info", "postal_code", postalCode, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch regional info")
		return
	}
	if ri == nil {
		c.String(http.StatusNotFound, "Regional info not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      ri.PostalCode,
		"DeleteURL": "/regional-infos/" + postalCode,
		"RowID":     "regional-info-row-" + postalCode,
	})
}
