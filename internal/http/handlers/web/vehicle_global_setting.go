package web

import (
	"log/slog"
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type VehicleGlobalSettingHandler struct {
	svc vehicle.VehicleGlobalSettingService
}

func NewVehicleGlobalSettingHandler(svc vehicle.VehicleGlobalSettingService) *VehicleGlobalSettingHandler {
	return &VehicleGlobalSettingHandler{svc: svc}
}

func (h *VehicleGlobalSettingHandler) List(c *gin.Context) {
	tz := c.DefaultQuery("tz", "UTC")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	settings, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		slog.Error("failed to list vehicle global settings", "error", err)
		c.String(http.StatusInternalServerError, "Failed to load vehicle service settings")
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

	c.HTML(http.StatusOK, "pages/vehicle_global_settings/index.html", gin.H{
		"title":           "Vehicle Service Settings",
		"settings":        settings,
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

func (h *VehicleGlobalSettingHandler) CreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/vehicle_global_settings/form.html", gin.H{
		"title":  "Add Vehicle Service Setting",
		"action": "/vehicle-global-settings",
	})
}

func (h *VehicleGlobalSettingHandler) Create(c *gin.Context) {
	key := c.PostForm("key")
	value := c.PostForm("value")
	remark := optionalPostForm(c, "remark")
	countryCode := optionalPostForm(c, "country_code")
	userEmail := c.GetString("user_email")

	_, err := h.svc.CreateVehicleGlobalSetting(key, value, remark, countryCode, userEmail)
	if err != nil {
		slog.Error("failed to create vehicle global setting", "key", key, "error", err)
		c.String(http.StatusInternalServerError, "Failed to create vehicle service setting")
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicle-global-settings")
}

func (h *VehicleGlobalSettingHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	setting, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch vehicle global setting", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to load vehicle service setting")
		return
	}
	if setting == nil {
		c.String(http.StatusNotFound, "Setting not found")
		return
	}

	c.HTML(http.StatusOK, "pages/vehicle_global_settings/form.html", gin.H{
		"title":   "Edit Vehicle Service Setting",
		"action":  "/vehicle-global-settings/" + idStr,
		"setting": setting,
	})
}

func (h *VehicleGlobalSettingHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	key := c.PostForm("key")
	value := c.PostForm("value")
	remark := optionalPostForm(c, "remark")
	countryCode := optionalPostForm(c, "country_code")
	userEmail := c.GetString("user_email")

	_, err := h.svc.UpdateVehicleGlobalSetting(id, key, value, remark, countryCode, userEmail)
	if err != nil {
		slog.Error("failed to update vehicle global setting", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to update vehicle service setting")
		return
	}

	c.Redirect(http.StatusSeeOther, "/vehicle-global-settings")
}

func (h *VehicleGlobalSettingHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteVehicleGlobalSetting(id); err != nil {
		slog.Error("failed to delete vehicle global setting", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to delete vehicle service setting")
		return
	}

	c.Status(http.StatusOK)
}

func (h *VehicleGlobalSettingHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tz := c.DefaultQuery("tz", "UTC")

	setting, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch vehicle global setting", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to load vehicle service setting")
		return
	}
	if setting == nil {
		c.String(http.StatusNotFound, "Setting not found")
		return
	}

	c.HTML(http.StatusOK, "pages/vehicle_global_settings/view.html", gin.H{
		"title":   "View Vehicle Service Setting",
		"setting": setting,
		"tz":      tz,
	})
}

func (h *VehicleGlobalSettingHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	setting, err := h.svc.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch vehicle global setting", "id", id, "error", err)
		c.String(http.StatusInternalServerError, "Failed to load vehicle service setting")
		return
	}
	if setting == nil {
		c.String(http.StatusNotFound, "Setting not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      setting.Key,
		"DeleteURL": "/vehicle-global-settings/" + idStr,
		"RowID":     "vehicle-global-setting-row-" + idStr,
	})
}

// optionalPostForm returns a pointer to the form value, or nil if empty.
func optionalPostForm(c *gin.Context, key string) *string {
	v := c.PostForm(key)
	if v == "" {
		return nil
	}
	return &v
}
