package api

import (
	"net/http"

	"golang-htmx-bulma/internal/user"
	"github.com/gin-gonic/gin"
)

type AuthAPIHandler struct {
	svc user.AuthService
}

// NewAuthAPIHandler creates a new AuthAPIHandler.
func NewAuthAPIHandler(svc user.AuthService) *AuthAPIHandler {
	return &AuthAPIHandler{svc: svc}
}

func (h *AuthAPIHandler) Login(c *gin.Context) {
	var payload struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, u, err := h.svc.Login(payload.Email, payload.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  u,
	})
}
