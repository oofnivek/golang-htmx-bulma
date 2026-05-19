package web

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"golang-htmx-bulma/internal/service"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) LoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/index.html", gin.H{
		"title":       "Sign In",
		"hideSidebar": true,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	token, _, err := h.authSvc.Login(email, password)
	if err != nil {
		// Log the failed login attempt using structured OpenTelemetry logging via slog.
		slog.Warn("Unsuccessful sign-in attempt",
			"email", email,
			"client_ip", c.ClientIP(),
			"error", err.Error(),
		)

		c.HTML(http.StatusUnauthorized, "pages/index.html", gin.H{
			"title":       "Sign In",
			"hideSidebar": true,
			"error":       err.Error(),
			"email":       email,
		})
		return
	}

	// Set HttpOnly cookie with the token.
	// Secure should ideally be set to true in production.
	c.SetCookie("jwt_token", token, 86400, "/", "", false, true)

	// Redirect to dashboard/users on successful login
	c.Redirect(http.StatusSeeOther, "/users")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear the cookie by setting max age to -1
	c.SetCookie("jwt_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/login")
}
