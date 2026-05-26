package routes

import (
	"golang-htmx-bulma/internal/http/handlers/api"
	"golang-htmx-bulma/internal/http/handlers/web"
	"golang-htmx-bulma/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all web and API routes. Handlers can be nil if running in an isolated service role.
func RegisterRoutes(
	r *gin.Engine,
	vcHandler *web.VehicleColorHandler,
	vmHandler *web.VehicleMakeHandler,
	vtHandler *web.VehicleTypeHandler,
	vsHandler *web.VehicleStatusHandler,
	ftHandler *web.FuelTypeHandler,
	vmdlHandler *web.VehicleModelHandler,
	vfHandler *web.VehicleFuelHandler,
	roleHandler *web.RoleHandler,
	userHandler *web.UserHandler,
	authHandler *web.AuthHandler,
	userAPI *api.UserAPIHandler,
	roleAPI *api.RoleAPIHandler,
	authAPI *api.AuthAPIHandler,
	vcAPI *api.VehicleColorAPIHandler,
	vmAPI *api.VehicleMakeAPIHandler,
	vsAPI *api.VehicleStatusAPIHandler,
	ftAPI *api.FuelTypeAPIHandler,
	vmdlAPI *api.VehicleModelAPIHandler,
	vfAPI *api.VehicleFuelAPIHandler,
	signingKey string,
) {
	// Static files
	r.Static("/static", "./static")

	// Auth & Web routes
	if authHandler != nil {
		r.GET("/", authHandler.LoginForm)
		r.GET("/login", authHandler.LoginForm)
		r.POST("/login", authHandler.Login)
		r.GET("/logout", authHandler.Logout)
	}

	// Protected Web Routes
	if vcHandler != nil || vmHandler != nil || vtHandler != nil || vsHandler != nil || ftHandler != nil || vmdlHandler != nil || vfHandler != nil || roleHandler != nil || userHandler != nil {
		protected := r.Group("/")
		protected.Use(middleware.AuthMiddleware(signingKey))
		{
			// Vehicle Colors
			if vcHandler != nil {
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
			}

			// Vehicle Makes
			if vmHandler != nil {
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
			}

			// Vehicle Types
			if vtHandler != nil {
				vtGroup := protected.Group("/vehicle-types")
				{
					vtGroup.GET("", vtHandler.List)
					vtGroup.GET("/new", vtHandler.CreateForm)
					vtGroup.POST("", vtHandler.Create)
					vtGroup.GET("/:id", vtHandler.View)
					vtGroup.GET("/:id/edit", vtHandler.EditForm)
					vtGroup.POST("/:id", vtHandler.Update)
					vtGroup.DELETE("/:id", vtHandler.Delete)
					vtGroup.GET("/:id/delete", vtHandler.DeleteConfirm)
				}
			}

			// Fuel Types
			if ftHandler != nil {
				ftGroup := protected.Group("/fuel-types")
				{
					ftGroup.GET("", ftHandler.List)
					ftGroup.GET("/new", ftHandler.CreateForm)
					ftGroup.POST("", ftHandler.Create)
					ftGroup.GET("/:id/view", ftHandler.View)
					ftGroup.GET("/:id/edit", ftHandler.EditForm)
					ftGroup.POST("/:id", ftHandler.Update)
					ftGroup.DELETE("/:id", ftHandler.Delete)
					ftGroup.GET("/:id/delete", ftHandler.DeleteConfirm)
				}
			}

			// Vehicle Models
			if vmdlHandler != nil {
				vmdlGroup := protected.Group("/vehicle-models")
				{
					vmdlGroup.GET("", vmdlHandler.List)
					vmdlGroup.GET("/new", vmdlHandler.CreateForm)
					vmdlGroup.POST("", vmdlHandler.Create)
					vmdlGroup.GET("/:id/edit", vmdlHandler.EditForm)
					vmdlGroup.POST("/:id", vmdlHandler.Update)
					vmdlGroup.DELETE("/:id", vmdlHandler.Delete)
					vmdlGroup.GET("/:id/delete", vmdlHandler.DeleteConfirm)
				}
			}

			// Vehicle Fuels
			if vfHandler != nil {
				vfGroup := protected.Group("/vehicle-fuels")
				{
					vfGroup.GET("", vfHandler.List)
					vfGroup.GET("/new", vfHandler.CreateForm)
					vfGroup.POST("", vfHandler.Create)
					vfGroup.GET("/:id/edit", vfHandler.EditForm)
					vfGroup.POST("/:id", vfHandler.Update)
					vfGroup.DELETE("/:id", vfHandler.Delete)
					vfGroup.GET("/:id/delete", vfHandler.DeleteConfirm)
				}
			}

			// Vehicle Statuses
			if vsHandler != nil {
				vsGroup := protected.Group("/vehicle-statuses")
				{
					vsGroup.GET("", vsHandler.List)
					vsGroup.GET("/new", vsHandler.CreateForm)
					vsGroup.POST("", vsHandler.Create)
					vsGroup.GET("/:id/view", vsHandler.View)
					vsGroup.GET("/:id/edit", vsHandler.EditForm)
					vsGroup.POST("/:id", vsHandler.Update)
					vsGroup.DELETE("/:id", vsHandler.Delete)
					vsGroup.GET("/:id/delete", vsHandler.DeleteConfirm)
				}
			}

			// Roles
			if roleHandler != nil {
				roleGroup := protected.Group("/roles")
				{
					roleGroup.GET("", roleHandler.List)
					roleGroup.GET("/new", roleHandler.CreateForm)
					roleGroup.POST("", roleHandler.Create)
					roleGroup.GET("/:id/view", roleHandler.View)
					roleGroup.GET("/:id/edit", roleHandler.EditForm)
					roleGroup.POST("/:id", roleHandler.Update)
					roleGroup.DELETE("/:id", roleHandler.Delete)
					roleGroup.GET("/:id/delete", roleHandler.DeleteConfirm)
				}
			}

			// Users
			if userHandler != nil {
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
		}
	}

	// API routes (REST and HTMX targets)
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/hello", api.HelloHandler)

		if userAPI != nil {
			apiGroup.GET("/users", userAPI.List)
			apiGroup.GET("/users/:email", userAPI.Get)
			apiGroup.POST("/users", userAPI.Create)
			apiGroup.PUT("/users/:email", userAPI.Update)
			apiGroup.DELETE("/users/:email", userAPI.Delete)
		}

		if roleAPI != nil {
			apiGroup.GET("/roles", roleAPI.List)
			apiGroup.GET("/roles/:id", roleAPI.Get)
			apiGroup.POST("/roles", roleAPI.Create)
			apiGroup.PUT("/roles/:id", roleAPI.Update)
			apiGroup.DELETE("/roles/:id", roleAPI.Delete)
		}

		if authAPI != nil {
			apiGroup.POST("/auth/login", authAPI.Login)
		}

		if vcAPI != nil {
			apiGroup.GET("/vehicle-colors/all", vcAPI.ListAll)
			apiGroup.GET("/vehicle-colors", vcAPI.List)
			apiGroup.GET("/vehicle-colors/:id", vcAPI.Get)
			apiGroup.POST("/vehicle-colors", vcAPI.Create)
			apiGroup.PUT("/vehicle-colors/:id", vcAPI.Update)
			apiGroup.DELETE("/vehicle-colors/:id", vcAPI.Delete)
		}

		if vmAPI != nil {
			apiGroup.GET("/vehicle-makes/all", vmAPI.ListAll)
			apiGroup.GET("/vehicle-makes", vmAPI.List)
			apiGroup.GET("/vehicle-makes/:id", vmAPI.Get)
			apiGroup.POST("/vehicle-makes", vmAPI.Create)
			apiGroup.PUT("/vehicle-makes/:id", vmAPI.Update)
			apiGroup.DELETE("/vehicle-makes/:id", vmAPI.Delete)
		}

		if vsAPI != nil {
			apiGroup.GET("/vehicle-statuses/all", vsAPI.ListAll)
			apiGroup.GET("/vehicle-statuses", vsAPI.List)
			apiGroup.GET("/vehicle-statuses/:id", vsAPI.Get)
			apiGroup.POST("/vehicle-statuses", vsAPI.Create)
			apiGroup.PUT("/vehicle-statuses/:id", vsAPI.Update)
			apiGroup.DELETE("/vehicle-statuses/:id", vsAPI.Delete)
		}

		if ftAPI != nil {
			apiGroup.GET("/fuel-types/all", ftAPI.ListAll)
			apiGroup.GET("/fuel-types", ftAPI.List)
			apiGroup.GET("/fuel-types/:id", ftAPI.Get)
			apiGroup.POST("/fuel-types", ftAPI.Create)
			apiGroup.PUT("/fuel-types/:id", ftAPI.Update)
			apiGroup.DELETE("/fuel-types/:id", ftAPI.Delete)
		}

		if vmdlAPI != nil {
			apiGroup.GET("/vehicle-models/all", vmdlAPI.ListAll)
			apiGroup.GET("/vehicle-models", vmdlAPI.List)
			apiGroup.GET("/vehicle-models/:id", vmdlAPI.Get)
			apiGroup.POST("/vehicle-models", vmdlAPI.Create)
			apiGroup.PUT("/vehicle-models/:id", vmdlAPI.Update)
			apiGroup.DELETE("/vehicle-models/:id", vmdlAPI.Delete)
		}

		if vfAPI != nil {
			apiGroup.GET("/vehicle-fuels/all", vfAPI.ListAll)
			apiGroup.GET("/vehicle-fuels", vfAPI.List)
			apiGroup.GET("/vehicle-fuels/:id", vfAPI.Get)
			apiGroup.POST("/vehicle-fuels", vfAPI.Create)
			apiGroup.PUT("/vehicle-fuels/:id", vfAPI.Update)
			apiGroup.DELETE("/vehicle-fuels/:id", vfAPI.Delete)
		}
	}
}
