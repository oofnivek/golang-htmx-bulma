package api

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type RegionalInfoAPIHandler struct {
	svc vehicle.RegionalInfoService
}

func NewRegionalInfoAPIHandler(svc vehicle.RegionalInfoService) *RegionalInfoAPIHandler {
	return &RegionalInfoAPIHandler{svc: svc}
}

func (h *RegionalInfoAPIHandler) ListAll(c *gin.Context) {
	items, err := h.svc.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"regional_infos": items})
}

func (h *RegionalInfoAPIHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "postal_code")
	sortOrder := c.DefaultQuery("sortOrder", "asc")

	items, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"regional_infos": items, "total": total})
}

func (h *RegionalInfoAPIHandler) Get(c *gin.Context) {
	postalCode := c.Param("id")

	ri, err := h.svc.FindByID(postalCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if ri == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "regional info not found"})
		return
	}
	c.JSON(http.StatusOK, ri)
}

func (h *RegionalInfoAPIHandler) Create(c *gin.Context) {
	var payload struct {
		PostalCode string `json:"postal_code" binding:"required"`
		Region     string `json:"region" binding:"required"`
		EstateID   int64  `json:"estate_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ri, err := h.svc.CreateRegionalInfo(payload.PostalCode, payload.Region, payload.EstateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ri)
}

func (h *RegionalInfoAPIHandler) Update(c *gin.Context) {
	postalCode := c.Param("id")

	var payload struct {
		Region   string `json:"region" binding:"required"`
		EstateID int64  `json:"estate_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ri, err := h.svc.UpdateRegionalInfo(postalCode, payload.Region, payload.EstateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if ri == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "regional info not found"})
		return
	}
	c.JSON(http.StatusOK, ri)
}

func (h *RegionalInfoAPIHandler) Delete(c *gin.Context) {
	postalCode := c.Param("id")

	if err := h.svc.DeleteRegionalInfo(postalCode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
