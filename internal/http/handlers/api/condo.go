package api

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/pkg/status"
	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type CondoAPIHandler struct{ svc vehicle.CondoService }

func NewCondoAPIHandler(svc vehicle.CondoService) *CondoAPIHandler {
	return &CondoAPIHandler{svc: svc}
}

func (h *CondoAPIHandler) ListAll(c *gin.Context) {
	condos, err := h.svc.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"condos": condos})
}

func (h *CondoAPIHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	condos, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"condos": condos,
		"total":  total,
	})
}

func (h *CondoAPIHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	condo, err := h.svc.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if condo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "condo not found"})
		return
	}
	c.JSON(http.StatusOK, condo)
}

func (h *CondoAPIHandler) Create(c *gin.Context) {
	var payload struct {
		Name       string        `json:"name" binding:"required"`
		Status     status.Status `json:"status"`
		McstNumber string        `json:"mcst_number" binding:"required"`
		McstEmail  string        `json:"mcst_email" binding:"required"`
		Address    string        `json:"address" binding:"required"`
		User       string        `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	condo, err := h.svc.CreateCondo(payload.Name, payload.Status, payload.McstNumber, payload.McstEmail, payload.Address, payload.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, condo)
}

func (h *CondoAPIHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var payload struct {
		Name       string        `json:"name" binding:"required"`
		Status     status.Status `json:"status"`
		McstNumber string        `json:"mcst_number" binding:"required"`
		McstEmail  string        `json:"mcst_email" binding:"required"`
		Address    string        `json:"address" binding:"required"`
		User       string        `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	condo, err := h.svc.UpdateCondo(id, payload.Name, payload.Status, payload.McstNumber, payload.McstEmail, payload.Address, payload.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if condo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "condo not found"})
		return
	}
	c.JSON(http.StatusOK, condo)
}

func (h *CondoAPIHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.svc.DeleteCondo(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
