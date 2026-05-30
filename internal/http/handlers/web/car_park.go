package web

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type CarParkHandler struct {
	svc    vehicle.CarParkService
	cpoSvc vehicle.CarParkOwnerService
}

func NewCarParkHandler(svc vehicle.CarParkService, cpoSvc vehicle.CarParkOwnerService) *CarParkHandler {
	return &CarParkHandler{svc: svc, cpoSvc: cpoSvc}
}

func (h *CarParkHandler) List(c *gin.Context) {
	tz := c.Query("tz")
	if tz == "" {
		tz = "UTC"
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	parks, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		slog.Error("failed to list car parks", "error", err)
		c.String(http.StatusInternalServerError, "Failed to list car parks")
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

	c.HTML(http.StatusOK, "pages/car_parks/index.html", gin.H{
		"title":           "Car Parks",
		"parks":           parks,
		"timezone":        tz,
		"timezones":       Timezones,
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

func (h *CarParkHandler) CreateForm(c *gin.Context) {
	owners, err := h.cpoSvc.ListAll()
	if err != nil {
		slog.Error("failed to list car park owners for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	c.HTML(http.StatusOK, "pages/car_parks/form.html", gin.H{
		"title":  "Add Car Park",
		"action": "/car-parks",
		"owners": owners,
	})
}

func (h *CarParkHandler) Create(c *gin.Context) {
	name := c.PostForm("name")
	description := nullableString(c.PostForm("description"))
	postalCode := c.PostForm("postal_code")
	address := c.PostForm("address")
	latitude, _ := strconv.ParseFloat(c.PostForm("latitude"), 64)
	longitude, _ := strconv.ParseFloat(c.PostForm("longitude"), 64)
	carParkOwnerID, _ := strconv.ParseInt(c.PostForm("car_park_owner_id"), 10, 64)
	activeFrom := parseDateTimeLocal(c.PostForm("active_from"))
	activeTo := parseDateTimeLocal(c.PostForm("active_to"))
	status := c.PostForm("status") == "on"
	userEmail := c.GetString("user_email")

	_, err := h.svc.CreateCarPark(name, description, postalCode, address, latitude, longitude, carParkOwnerID, activeFrom, activeTo, status, userEmail)
	if err != nil {
		slog.Error("failed to create car park", "error", err)
		c.String(http.StatusInternalServerError, "Failed to create car park")
		return
	}

	c.Redirect(http.StatusSeeOther, "/car-parks")
}

func (h *CarParkHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	cp, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch car park", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch car park")
		return
	}
	if cp == nil {
		c.String(http.StatusNotFound, "Car park not found")
		return
	}

	owners, err := h.cpoSvc.ListAll()
	if err != nil {
		slog.Error("failed to list car park owners for form", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load form data")
		return
	}

	c.HTML(http.StatusOK, "pages/car_parks/form.html", gin.H{
		"title":          "Edit Car Park",
		"action":         "/car-parks/" + idStr,
		"park":           cp,
		"owners":         owners,
		"activeFromStr":  formatDateTimeLocal(cp.ActiveFrom),
		"activeToStr":    formatDateTimeLocal(cp.ActiveTo),
	})
}

func (h *CarParkHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	name := c.PostForm("name")
	description := nullableString(c.PostForm("description"))
	postalCode := c.PostForm("postal_code")
	address := c.PostForm("address")
	latitude, _ := strconv.ParseFloat(c.PostForm("latitude"), 64)
	longitude, _ := strconv.ParseFloat(c.PostForm("longitude"), 64)
	carParkOwnerID, _ := strconv.ParseInt(c.PostForm("car_park_owner_id"), 10, 64)
	activeFrom := parseDateTimeLocal(c.PostForm("active_from"))
	activeTo := parseDateTimeLocal(c.PostForm("active_to"))
	status := c.PostForm("status") == "on"
	userEmail := c.GetString("user_email")

	_, err := h.svc.UpdateCarPark(id, name, description, postalCode, address, latitude, longitude, carParkOwnerID, activeFrom, activeTo, status, userEmail)
	if err != nil {
		slog.Error("failed to update car park", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to update car park")
		return
	}

	c.Redirect(http.StatusSeeOther, "/car-parks")
}

func (h *CarParkHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteCarPark(id); err != nil {
		slog.Error("failed to delete car park", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete car park")
		return
	}

	c.Status(http.StatusOK)
}

func (h *CarParkHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tz := c.DefaultQuery("tz", "UTC")

	cp, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch car park", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch car park")
		return
	}
	if cp == nil {
		c.String(http.StatusNotFound, "Car park not found")
		return
	}

	c.HTML(http.StatusOK, "pages/car_parks/view.html", gin.H{
		"title": "View Car Park",
		"park":  cp,
		"tz":    tz,
	})
}

func (h *CarParkHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	cp, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch car park", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to fetch car park")
		return
	}
	if cp == nil {
		c.String(http.StatusNotFound, "Car park not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      cp.Name,
		"DeleteURL": "/car-parks/" + idStr,
		"RowID":     "car-park-row-" + idStr,
	})
}

func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func parseDateTimeLocal(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02T15:04", s)
	if err != nil {
		return nil
	}
	utc := t.UTC()
	return &utc
}

func formatDateTimeLocal(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format("2006-01-02T15:04")
}
