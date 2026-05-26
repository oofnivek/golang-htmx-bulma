package web

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type FuelCompanyHandler struct {
	svc vehicle.FuelCompanyService
}

func NewFuelCompanyHandler(svc vehicle.FuelCompanyService) *FuelCompanyHandler {
	return &FuelCompanyHandler{svc: svc}
}

func (h *FuelCompanyHandler) List(c *gin.Context) {
	tz := c.Query("tz")
	if tz == "" {
		tz = "UTC"
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	companies, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	c.HTML(http.StatusOK, "pages/fuel_companies/index.html", gin.H{
		"title":           "Fuel Companies",
		"companies":       companies,
		"timezone":        tz,
		"timezones":       Timezones,
		"currentPage":     page,
		"pageSize":        pageSize,
		"pageSizeOptions": []int{5, 10, 20, 50},
		"totalPages":      totalPages,
		"total":           total,
		"sortBy":          sortBy,
		"sortOrder":       sortOrder,
	})
}

func (h *FuelCompanyHandler) CreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/fuel_companies/form.html", gin.H{
		"title":  "Add Fuel Company",
		"action": "/fuel-companies",
	})
}

func (h *FuelCompanyHandler) Create(c *gin.Context) {
	name := c.PostForm("name")
	status := c.PostForm("status") == "on"

	userEmail := c.GetString("user_email")
	_, err := h.svc.CreateFuelCompany(name, status, userEmail)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/fuel-companies")
}

func (h *FuelCompanyHandler) EditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	company, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if company == nil {
		c.String(http.StatusNotFound, "Fuel company not found")
		return
	}

	c.HTML(http.StatusOK, "pages/fuel_companies/form.html", gin.H{
		"title":   "Edit Fuel Company",
		"action":  "/fuel-companies/" + idStr,
		"company": company,
	})
}

func (h *FuelCompanyHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	name := c.PostForm("name")
	status := c.PostForm("status") == "on"

	userEmail := c.GetString("user_email")
	_, err := h.svc.UpdateFuelCompany(id, name, status, userEmail)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/fuel-companies")
}

func (h *FuelCompanyHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.svc.DeleteFuelCompany(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *FuelCompanyHandler) View(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tz := c.DefaultQuery("tz", "UTC")

	company, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if company == nil {
		c.String(http.StatusNotFound, "Fuel company not found")
		return
	}

	c.HTML(http.StatusOK, "pages/fuel_companies/view.html", gin.H{
		"title":   "View Fuel Company",
		"company": company,
		"tz":      tz,
	})
}

func (h *FuelCompanyHandler) DeleteConfirm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	company, err := h.svc.FindByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if company == nil {
		c.String(http.StatusNotFound, "Fuel company not found")
		return
	}

	c.HTML(http.StatusOK, "partials/modals/delete_confirm.html", gin.H{
		"Name":      company.Name,
		"DeleteURL": "/fuel-companies/" + idStr,
		"RowID":     "company-row-" + idStr,
	})
}
