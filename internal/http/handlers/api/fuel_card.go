package api

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/vehicle"

	"github.com/gin-gonic/gin"
)

type FuelCardAPIHandler struct {
	svc vehicle.FuelCardService
}

func NewFuelCardAPIHandler(svc vehicle.FuelCardService) *FuelCardAPIHandler {
	return &FuelCardAPIHandler{svc: svc}
}

func (h *FuelCardAPIHandler) ListAll(c *gin.Context) {
	cards, err := h.svc.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cards": cards})
}

func (h *FuelCardAPIHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	cards, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cards": cards,
		"total": total,
	})
}

func (h *FuelCardAPIHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	card, err := h.svc.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if card == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "fuel card not found"})
		return
	}
	c.JSON(http.StatusOK, card)
}

func (h *FuelCardAPIHandler) Create(c *gin.Context) {
	var payload struct {
		CardNo        string `json:"card_no" binding:"required"`
		FuelCompanyID int64  `json:"fuel_company_id" binding:"required"`
		PinNumber     string `json:"pin_number" binding:"required"`
		VehicleID     *int64 `json:"vehicle_id"`
		Status        bool   `json:"status"`
		User          string `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card, err := h.svc.CreateFuelCard(payload.CardNo, payload.FuelCompanyID, payload.PinNumber, payload.VehicleID, payload.Status, payload.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, card)
}

func (h *FuelCardAPIHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var payload struct {
		CardNo        string `json:"card_no" binding:"required"`
		FuelCompanyID int64  `json:"fuel_company_id" binding:"required"`
		PinNumber     string `json:"pin_number" binding:"required"`
		VehicleID     *int64 `json:"vehicle_id"`
		Status        bool   `json:"status"`
		User          string `json:"user" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card, err := h.svc.UpdateFuelCard(id, payload.CardNo, payload.FuelCompanyID, payload.PinNumber, payload.VehicleID, payload.Status, payload.User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if card == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "fuel card not found"})
		return
	}
	c.JSON(http.StatusOK, card)
}

func (h *FuelCardAPIHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	if err := h.svc.DeleteFuelCard(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
