package web

import (
	"log/slog"
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type CarAssetOwnerHandler struct {
	svc vehicle.CarAssetOwnerService
}

func NewCarAssetOwnerHandler(svc vehicle.CarAssetOwnerService) *CarAssetOwnerHandler {
	return &CarAssetOwnerHandler{svc: svc}
}

func (h *CarAssetOwnerHandler) List(c *gin.Context) {
	tz := c.Query("tz")
	if tz == "" {
		tz = "UTC"
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	owners, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		slog.Error("failed to list car asset owners", "error", err)
		c.String(http.StatusInternalServerError, "Failed to list car asset owners")
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

	c.HTML(http.StatusOK, "pages/car_asset_owners/index.html", gin.H{
		"title":           "Car Asset Owners",
		"owners":          owners,
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
	})
}

func (h *CarAssetOwnerHandler) CreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/car_asset_owners/form.html", gin.H{
		"title":  "Add Car Asset Owner",
		"action": "/car-asset-owners",
	})
}

func (h *CarAssetOwnerHandler) Create(c *gin.Context) {
	name := c.PostForm("name")
	status := c.PostForm("status") == "on"
	userEmail := c.GetString("user_email")

	_, err := h.svc.CreateOwner(name, status, userEmail)
	if err != nil {
		slog.Error("failed to create car asset owner", "error", err)
		c.String(http.StatusInternalServerError, "Failed to create car asset owner")
		return
	}

	c.Redirect(http.StatusSeeOther, "/car-asset-owners")
}

func (h *CarAssetOwnerHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	owner, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch car asset owner", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch car asset owner")
		return
	}
	if owner == nil {
		c.String(http.StatusNotFound, "Car asset owner not found")
		return
	}

	c.HTML(http.StatusOK, "pages/car_asset_owners/form.html", gin.H{
		"title":  "Edit Car Asset Owner",
		"action": "/car-asset-owners/" + idStr,
		"owner":  owner,
	})
}

func (h *CarAssetOwnerHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	name := c.PostForm("name")
	status := c.PostForm("status") == "on"
	userEmail := c.GetString("user_email")

	_, err := h.svc.UpdateOwner(id, name, status, userEmail)
	if err != nil {
		slog.Error("failed to update car asset owner", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to update car asset owner")
		return
	}

	c.Redirect(http.StatusSeeOther, "/car-asset-owners")
}

func (h *CarAssetOwnerHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteOwner(id); err != nil {
		slog.Error("failed to delete car asset owner", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete car asset owner")
		return
	}

	c.Status(http.StatusOK)
}

func (h *CarAssetOwnerHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tz := c.DefaultQuery("tz", "UTC")

	owner, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch car asset owner", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch car asset owner")
		return
	}
	if owner == nil {
		c.String(http.StatusNotFound, "Car asset owner not found")
		return
	}

	c.HTML(http.StatusOK, "pages/car_asset_owners/view.html", gin.H{
		"title": "View Car Asset Owner",
		"owner": owner,
		"tz":    tz,
	})
}

func (h *CarAssetOwnerHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	owner, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch car asset owner", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch car asset owner")
		return
	}
	if owner == nil {
		c.String(http.StatusNotFound, "Car asset owner not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      owner.Name,
		"DeleteURL": "/car-asset-owners/" + idStr,
		"RowID":     "owner-row-" + idStr,
	})
}
