package api

import (
	"net/http"
	"strconv"
	"time"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type CarParkAPIHandler struct {
	svc vehicle.CarParkService
}

func NewCarParkAPIHandler(svc vehicle.CarParkService) *CarParkAPIHandler {
	return &CarParkAPIHandler{svc: svc}
}

func (h *CarParkAPIHandler) ListAll(c *gin.Context) {
	parks, err := h.svc.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"parks": parks})
}

func (h *CarParkAPIHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	parks, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"parks": parks,
		"total": total,
	})
}

func (h *CarParkAPIHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	cp, err := h.svc.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car park not found"})
		return
	}
	c.JSON(http.StatusOK, cp)
}

func (h *CarParkAPIHandler) Create(c *gin.Context) {
	var payload struct {
		Name           string     `json:"name" binding:"required"`
		Description    *string    `json:"description"`
		PostalCode     string     `json:"postal_code" binding:"required"`
		Address        string     `json:"address" binding:"required"`
		Latitude       float64    `json:"latitude"`
		Longitude      float64    `json:"longitude"`
		CarParkOwnerID int64      `json:"car_park_owner_id"`
		ActiveFrom     *time.Time `json:"active_from"`
		ActiveTo       *time.Time `json:"active_to"`
		Status         bool       `json:"status"`
		User           string     `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cp, err := h.svc.CreateCarPark(
		payload.Name, payload.Description, payload.PostalCode, payload.Address,
		payload.Latitude, payload.Longitude, payload.CarParkOwnerID,
		payload.ActiveFrom, payload.ActiveTo, payload.Status, payload.User,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cp)
}

func (h *CarParkAPIHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var payload struct {
		Name           string     `json:"name" binding:"required"`
		Description    *string    `json:"description"`
		PostalCode     string     `json:"postal_code" binding:"required"`
		Address        string     `json:"address" binding:"required"`
		Latitude       float64    `json:"latitude"`
		Longitude      float64    `json:"longitude"`
		CarParkOwnerID int64      `json:"car_park_owner_id"`
		ActiveFrom     *time.Time `json:"active_from"`
		ActiveTo       *time.Time `json:"active_to"`
		Status         bool       `json:"status"`
		User           string     `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cp, err := h.svc.UpdateCarPark(
		id, payload.Name, payload.Description, payload.PostalCode, payload.Address,
		payload.Latitude, payload.Longitude, payload.CarParkOwnerID,
		payload.ActiveFrom, payload.ActiveTo, payload.Status, payload.User,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car park not found"})
		return
	}
	c.JSON(http.StatusOK, cp)
}

func (h *CarParkAPIHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.svc.DeleteCarPark(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
