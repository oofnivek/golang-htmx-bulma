package api

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/user"
	"github.com/gin-gonic/gin"
)

type RoleAPIHandler struct {
	svc user.RoleService
}

// NewRoleAPIHandler creates a new RoleAPIHandler.
func NewRoleAPIHandler(svc user.RoleService) *RoleAPIHandler {
	return &RoleAPIHandler{svc: svc}
}

func (h *RoleAPIHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "asc")
	search := c.Query("search")

	roles, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": roles,
		"total": total,
	})
}

func (h *RoleAPIHandler) Get(c *gin.Context) {
	id := c.Param("id")
	r, err := h.svc.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if r == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}
	c.JSON(http.StatusOK, r)
}

func (h *RoleAPIHandler) Create(c *gin.Context) {
	var payload struct {
		ID   string `json:"id" binding:"required"`
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.svc.CreateRole(payload.ID, payload.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, role)
}

func (h *RoleAPIHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var payload struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.svc.UpdateRole(id, payload.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if role == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}
	c.JSON(http.StatusOK, role)
}

func (h *RoleAPIHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteRole(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
