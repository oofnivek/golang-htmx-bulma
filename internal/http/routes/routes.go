package routes

import (
	"golang-htmx-bulma/internal/http/handlers/api"
	"golang-htmx-bulma/internal/http/handlers/web"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, vcHandler *web.VehicleColorHandler, vmHandler *web.VehicleMakeHandler) {
	// Static files
	r.Static("/static", "./static")

	// Web routes
	r.GET("/", web.HomeHandler)

	// Vehicle Colors
	vcGroup := r.Group("/vehicle-colors")
	{
		vcGroup.GET("", vcHandler.List)
		vcGroup.GET("/new", vcHandler.CreateForm)
		vcGroup.POST("", vcHandler.Create)
		vcGroup.GET("/:id/edit", vcHandler.EditForm)
		vcGroup.POST("/:id", vcHandler.Update)
		vcGroup.DELETE("/:id", vcHandler.Delete)
		vcGroup.GET("/:id/delete", vcHandler.DeleteConfirm)
	}

	// Vehicle Makes
	vmGroup := r.Group("/vehicle-makes")
	{
		vmGroup.GET("", vmHandler.List)
		vmGroup.GET("/new", vmHandler.CreateForm)
		vmGroup.POST("", vmHandler.Create)
		vmGroup.GET("/:id/edit", vmHandler.EditForm)
		vmGroup.POST("/:id", vmHandler.Update)
		vmGroup.DELETE("/:id", vmHandler.Delete)
		vmGroup.GET("/:id/delete", vmHandler.DeleteConfirm)
	}

	// API routes (HTMX targets)
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/hello", api.HelloHandler)
	}
}
