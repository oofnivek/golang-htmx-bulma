package web

import (
	"log/slog"
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type CarParkOwnerHandler struct {
	svc vehicle.CarParkOwnerService
}

func NewCarParkOwnerHandler(svc vehicle.CarParkOwnerService) *CarParkOwnerHandler {
	return &CarParkOwnerHandler{svc: svc}
}

func (h *CarParkOwnerHandler) List(c *gin.Context) {
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
		slog.Error("failed to list car park owners", "error", err)
		c.String(http.StatusInternalServerError, "Failed to list car park owners")
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

	c.HTML(http.StatusOK, "pages/car_park_owners/index.html", gin.H{
		"title":           "Car Park Owners",
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

func (h *CarParkOwnerHandler) CreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/car_park_owners/form.html", gin.H{
		"title":  "Add Car Park Owner",
		"action": "/car-park-owners",
	})
}

func (h *CarParkOwnerHandler) Create(c *gin.Context) {
	name := c.PostForm("name")
	status := c.PostForm("status") == "on"
	userEmail := c.GetString("user_email")

	_, err := h.svc.CreateOwner(name, status, userEmail)
	if err != nil {
		slog.Error("failed to create car park owner", "error", err)
		c.String(http.StatusInternalServerError, "Failed to create car park owner")
		return
	}

	c.Redirect(http.StatusSeeOther, "/car-park-owners")
}

func (h *CarParkOwnerHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	owner, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch car park owner", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch car park owner")
		return
	}
	if owner == nil {
		c.String(http.StatusNotFound, "Car park owner not found")
		return
	}

	c.HTML(http.StatusOK, "pages/car_park_owners/form.html", gin.H{
		"title":  "Edit Car Park Owner",
		"action": "/car-park-owners/" + idStr,
		"owner":  owner,
	})
}

func (h *CarParkOwnerHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	name := c.PostForm("name")
	status := c.PostForm("status") == "on"
	userEmail := c.GetString("user_email")

	_, err := h.svc.UpdateOwner(id, name, status, userEmail)
	if err != nil {
		slog.Error("failed to update car park owner", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to update car park owner")
		return
	}

	c.Redirect(http.StatusSeeOther, "/car-park-owners")
}

func (h *CarParkOwnerHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteOwner(id); err != nil {
		slog.Error("failed to delete car park owner", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete car park owner")
		return
	}

	c.Status(http.StatusOK)
}

func (h *CarParkOwnerHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tz := c.DefaultQuery("tz", "UTC")

	owner, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch car park owner", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch car park owner")
		return
	}
	if owner == nil {
		c.String(http.StatusNotFound, "Car park owner not found")
		return
	}

	c.HTML(http.StatusOK, "pages/car_park_owners/view.html", gin.H{
		"title": "View Car Park Owner",
		"owner": owner,
		"tz":    tz,
	})
}

func (h *CarParkOwnerHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	owner, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch car park owner", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch car park owner")
		return
	}
	if owner == nil {
		c.String(http.StatusNotFound, "Car park owner not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      owner.Name,
		"DeleteURL": "/car-park-owners/" + idStr,
		"RowID":     "park-owner-row-" + idStr,
	})
}
