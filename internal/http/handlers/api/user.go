package api

import (
	"net/http"
	"strconv"

	"golang-htmx-bulma/internal/user"
	"github.com/gin-gonic/gin"
)

type UserAPIHandler struct {
	svc user.UserService
}

// NewUserAPIHandler creates a new UserAPIHandler.
func NewUserAPIHandler(svc user.UserService) *UserAPIHandler {
	return &UserAPIHandler{svc: svc}
}

func (h *UserAPIHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "email")
	sortOrder := c.DefaultQuery("sortOrder", "asc")
	users, total, err := h.svc.ListPaged(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
	})
}

func (h *UserAPIHandler) Get(c *gin.Context) {
	email := c.Param("email")
	u, err := h.svc.GetByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if u == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserAPIHandler) Create(c *gin.Context) {
	var u user.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.CreateUser(&u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, u)
}

func (h *UserAPIHandler) Update(c *gin.Context) {
	email := c.Param("email")
	var u user.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u.Email = email

	if err := h.svc.UpdateUser(&u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserAPIHandler) Delete(c *gin.Context) {
	email := c.Param("email")
	if err := h.svc.DeleteUser(email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
