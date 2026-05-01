package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/index.html", gin.H{
		"title": "Home Page",
	})
}
