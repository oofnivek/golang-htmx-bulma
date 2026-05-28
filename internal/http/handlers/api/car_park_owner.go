package api

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type CarParkOwnerAPIHandler struct {
	svc vehicle.CarParkOwnerService
}

func NewCarParkOwnerAPIHandler(svc vehicle.CarParkOwnerService) *CarParkOwnerAPIHandler {
	return &CarParkOwnerAPIHandler{svc: svc}
}

func (h *CarParkOwnerAPIHandler) ListAll(c *gin.Context) {
	owners, err := h.svc.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"owners": owners})
}

func (h *CarParkOwnerAPIHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	owners, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"owners": owners,
		"total":  total,
	})
}

func (h *CarParkOwnerAPIHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	owner, err := h.svc.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if owner == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car park owner not found"})
		return
	}
	c.JSON(http.StatusOK, owner)
}

func (h *CarParkOwnerAPIHandler) Create(c *gin.Context) {
	var payload struct {
		Name   string `json:"name" binding:"required"`
		Status bool   `json:"status"`
		User   string `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	owner, err := h.svc.CreateOwner(payload.Name, payload.Status, payload.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, owner)
}

func (h *CarParkOwnerAPIHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var payload struct {
		Name   string `json:"name" binding:"required"`
		Status bool   `json:"status"`
		User   string `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	owner, err := h.svc.UpdateOwner(id, payload.Name, payload.Status, payload.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if owner == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car park owner not found"})
		return
	}
	c.JSON(http.StatusOK, owner)
}

func (h *CarParkOwnerAPIHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.svc.DeleteOwner(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
