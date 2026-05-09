package routes

import (
	"golang-htmx-bulma/internal/http/handlers/api"
	"golang-htmx-bulma/internal/http/handlers/web"
	"golang-htmx-bulma/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, vcHandler *web.VehicleColorHandler, vmHandler *web.VehicleMakeHandler, roleHandler *web.RoleHandler, userHandler *web.UserHandler, authHandler *web.AuthHandler, signingKey string) {
	// Static files
	r.Static("/static", "./static")

	// Auth & Web routes
	r.GET("/", authHandler.LoginForm)
	r.GET("/login", authHandler.LoginForm)
	r.POST("/login", authHandler.Login)
	r.GET("/logout", authHandler.Logout)

	// Protected Web Routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(signingKey))
	{
		// Vehicle Colors
		vcGroup := protected.Group("/vehicle-colors")
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
		vmGroup := protected.Group("/vehicle-makes")
		{
			vmGroup.GET("", vmHandler.List)
			vmGroup.GET("/new", vmHandler.CreateForm)
			vmGroup.POST("", vmHandler.Create)
			vmGroup.GET("/:id/edit", vmHandler.EditForm)
			vmGroup.POST("/:id", vmHandler.Update)
			vmGroup.DELETE("/:id", vmHandler.Delete)
			vmGroup.GET("/:id/delete", vmHandler.DeleteConfirm)
		}

		// Roles
		roleGroup := protected.Group("/roles")
		{
			roleGroup.GET("", roleHandler.List)
			roleGroup.GET("/new", roleHandler.CreateForm)
			roleGroup.POST("", roleHandler.Create)
			roleGroup.GET("/:id/edit", roleHandler.EditForm)
			roleGroup.POST("/:id", roleHandler.Update)
			roleGroup.DELETE("/:id", roleHandler.Delete)
			roleGroup.GET("/:id/delete", roleHandler.DeleteConfirm)
		}

		// Users
		userGroup := protected.Group("/users")
		{
			userGroup.GET("", userHandler.List)
			userGroup.GET("/new", userHandler.CreateForm)
			userGroup.GET("/password-fields", userHandler.PasswordFields)
			userGroup.POST("", userHandler.Create)
			userGroup.GET("/:email/edit", userHandler.EditForm)
			userGroup.GET("/:email/view", userHandler.View)
			userGroup.POST("/:email", userHandler.Update)
			userGroup.DELETE("/:email", userHandler.Delete)
			userGroup.GET("/:email/delete", userHandler.DeleteConfirm)
		}
	}

	// API routes (HTMX targets)
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/hello", api.HelloHandler)
	}
}
