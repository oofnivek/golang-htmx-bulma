package api

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type CarParkLotAPIHandler struct {
	svc vehicle.CarParkLotService
}

func NewCarParkLotAPIHandler(svc vehicle.CarParkLotService) *CarParkLotAPIHandler {
	return &CarParkLotAPIHandler{svc: svc}
}

func (h *CarParkLotAPIHandler) ListAll(c *gin.Context) {
	lots, err := h.svc.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"lots": lots})
}

func (h *CarParkLotAPIHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	lots, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"lots":  lots,
		"total": total,
	})
}

func (h *CarParkLotAPIHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	l, err := h.svc.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if l == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car park lot not found"})
		return
	}
	c.JSON(http.StatusOK, l)
}

func (h *CarParkLotAPIHandler) Create(c *gin.Context) {
	var payload struct {
		CarParkID int64  `json:"car_park_id" binding:"required"`
		LotNumber string `json:"lot_number" binding:"required"`
		Level     string `json:"level" binding:"required"`
		Status    bool   `json:"status"`
		User      string `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	l, err := h.svc.CreateCarParkLot(payload.CarParkID, payload.LotNumber, payload.Level, payload.Status, payload.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, l)
}

func (h *CarParkLotAPIHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var payload struct {
		CarParkID int64  `json:"car_park_id" binding:"required"`
		LotNumber string `json:"lot_number" binding:"required"`
		Level     string `json:"level" binding:"required"`
		Status    bool   `json:"status"`
		User      string `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	l, err := h.svc.UpdateCarParkLot(id, payload.CarParkID, payload.LotNumber, payload.Level, payload.Status, payload.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if l == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car park lot not found"})
		return
	}
	c.JSON(http.StatusOK, l)
}

func (h *CarParkLotAPIHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.svc.DeleteCarParkLot(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
