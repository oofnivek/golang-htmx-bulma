package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloHandler(c *gin.Context) {
	// Return a small HTML fragment for HTMX
	c.String(http.StatusOK, `<div class="notification is-success">Hello from the server! HTMX is working.</div>`)
}
