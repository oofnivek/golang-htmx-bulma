package web

import (
    "net/http"
    "strconv"

    "golang-htmx-bulma/internal/vehicle"

    "github.com/gin-gonic/gin"
)

type VehicleTypeHandler struct {
    svc vehicle.VehicleTypeService
}

func NewVehicleTypeHandler(svc vehicle.VehicleTypeService) *VehicleTypeHandler {
    return &VehicleTypeHandler{svc: svc}
}

func RegisterVehicleTypeRoutes(r *gin.RouterGroup, svc vehicle.VehicleTypeService) {
    h := NewVehicleTypeHandler(svc)
    r.GET("/vehicle-types", h.List)
    r.GET("/vehicle-types/new", h.CreateForm)
    r.POST("/vehicle-types", h.Create)
    r.GET("/vehicle-types/:id", h.View)
    r.GET("/vehicle-types/:id/edit", h.EditForm)
    r.POST("/vehicle-types/:id", h.Update)
    r.POST("/vehicle-types/:id/delete", h.Delete)
    r.GET("/vehicle-types/:id/delete", h.DeleteConfirm)
}

func (h *VehicleTypeHandler) List(c *gin.Context) {
    tz := c.Query("tz")
    if tz == "" { tz = "UTC" }
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
    sortBy := c.DefaultQuery("sortBy", "id")
    sortOrder := c.DefaultQuery("sortOrder", "desc")
    types, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
    if err != nil { c.String(http.StatusInternalServerError, err.Error()); return }
    timezones := Timezones
    totalPages := (total + pageSize - 1) / pageSize
    c.HTML(http.StatusOK, "pages/vehicle_type/index.html", gin.H{
        "title":          "Vehicle Types",
        "types":          types,
        "timezone":       tz,
        "timezones":      timezones,
        "currentPage":    page,
        "pageSize":       pageSize,
        "pageSizeOptions": []int{5,10,20,50},
        "totalPages":     totalPages,
        "total":          total,
        "sortBy":         sortBy,
        "sortOrder":      sortOrder,
    })
}

func (h *VehicleTypeHandler) CreateForm(c *gin.Context) {
    c.HTML(http.StatusOK, "pages/vehicle_type/form.html", gin.H{"title": "Add Vehicle Type", "action": "/vehicle-types"})
}

func (h *VehicleTypeHandler) Create(c *gin.Context) {
    name := c.PostForm("name")
    statusStr := c.PostForm("status")
    status := statusStr == "on"
    userEmail := c.GetString("user_email")
    _, err := h.svc.Create(name, status, userEmail)
    if err != nil { c.String(http.StatusInternalServerError, err.Error()); return }
    c.Redirect(http.StatusSeeOther, "/vehicle-types")
}

func (h *VehicleTypeHandler) EditForm(c *gin.Context) {
    idStr := c.Param("id")
    id, _ := strconv.ParseInt(idStr, 10, 64)
    vt, err := h.svc.FindByID(id)
    if err != nil { c.String(http.StatusInternalServerError, err.Error()); return }
    if vt == nil { c.String(http.StatusNotFound, "Vehicle type not found"); return }
    c.HTML(http.StatusOK, "pages/vehicle_type/form.html", gin.H{"title": "Edit Vehicle Type", "action": "/vehicle-types/" + idStr, "type": vt})
}

func (h *VehicleTypeHandler) Update(c *gin.Context) {
    idStr := c.Param("id")
    id, _ := strconv.ParseInt(idStr, 10, 64)
    name := c.PostForm("name")
    statusStr := c.PostForm("status")
    status := statusStr == "on"
    userEmail := c.GetString("user_email")
    _, err := h.svc.Update(id, name, status, userEmail)
    if err != nil { c.String(http.StatusInternalServerError, err.Error()); return }
    c.Redirect(http.StatusSeeOther, "/vehicle-types")
}

func (h *VehicleTypeHandler) Delete(c *gin.Context) {
    idStr := c.Param("id")
    id, _ := strconv.ParseInt(idStr, 10, 64)
    err := h.svc.Delete(id)
    if err != nil { c.String(http.StatusInternalServerError, err.Error()); return }
    c.Status(http.StatusOK)
}

func (h *VehicleTypeHandler) DeleteConfirm(c *gin.Context) {
    idStr := c.Param("id")
    id, _ := strconv.ParseInt(idStr, 10, 64)
    vt, err := h.svc.FindByID(id)
    if err != nil { c.String(http.StatusInternalServerError, err.Error()); return }
    if vt == nil { c.String(http.StatusNotFound, "Vehicle type not found"); return }
    c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{"Name": vt.Name, "DeleteURL": "/vehicle-types/" + idStr, "RowID": "type-row-" + idStr})
}

func (h *VehicleTypeHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tz := c.DefaultQuery("tz", "UTC")

	vt, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if vt == nil {
		c.String(http.StatusNotFound, "Vehicle type not found")
		return
	}

	c.HTML(http.StatusOK, "pages/vehicle_type/view.html", gin.H{
		"title": "View Vehicle Type",
		"type":  vt,
		"tz":    tz,
	})
}
