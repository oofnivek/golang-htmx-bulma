package main

import (
	"context"
	"log"
	"net/http"

	"golang-htmx-bulma/internal/config"
	"golang-htmx-bulma/internal/db"
	"golang-htmx-bulma/internal/http/handlers/web"
	"golang-htmx-bulma/internal/http/routes"
	"golang-htmx-bulma/internal/repository"
	"golang-htmx-bulma/internal/service"
	"golang-htmx-bulma/internal/telemetry"
	"golang-htmx-bulma/internal/vehicle"
	"golang-htmx-bulma/internal/view"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.Load()

	// Initialize OpenTelemetry Logging
	ctx := context.Background()
	logProvider, err := telemetry.InitLogger(ctx, "fleet-management-system", cfg.AppEnv)
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry logging: %v\n", err)
	}
	defer func() {
		if err := logProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down OpenTelemetry logger provider: %v\n", err)
		}
	}()

	// Initialize Vehicle DB
	vehicleDatabase, err := db.InitDB(cfg.VehicleDBDSN)
	if err != nil {
		log.Fatalf("Failed to connect to Vehicle MySQL: %v\n", err)
	}
	defer vehicleDatabase.Close()

	// Initialize User DB
	userDatabase, err := db.InitDB(cfg.UserDBDSN)
	if err != nil {
		log.Fatalf("Failed to connect to User MySQL: %v\n", err)
	}
	defer userDatabase.Close()

	// Initialize Gin
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Load templates
	r.HTMLRender = view.NewRenderer("templates")

	// Register routes
	vcRepo := vehicle.NewVehicleColorRepository(vehicleDatabase)
	vcSvc := vehicle.NewVehicleColorService(vcRepo)
	vcHandler := web.NewVehicleColorHandler(vcSvc)

	vmRepo := vehicle.NewVehicleMakeRepository(vehicleDatabase)
	vmSvc := vehicle.NewVehicleMakeService(vmRepo)
	vmHandler := web.NewVehicleMakeHandler(vmSvc)

	roleRepo := repository.NewRoleRepository(userDatabase)
	roleSvc := service.NewRoleService(roleRepo)
	roleHandler := web.NewRoleHandler(roleSvc)

	userRepo := repository.NewUserRepository(userDatabase)
	userSvc := service.NewUserService(userRepo)
	userHandler := web.NewUserHandler(userSvc, roleSvc)

	if cfg.JWTSigningKey == "" {
		log.Fatal("JWT_SIGNING_KEY environment variable is not set")
	}
	authSvc := service.NewAuthService(userSvc, cfg.JWTSigningKey)
	authHandler := web.NewAuthHandler(authSvc)

	routes.RegisterRoutes(r, vcHandler, vmHandler, roleHandler, userHandler, authHandler, cfg.JWTSigningKey)

	// Start server
	log.Printf("Server starting on port %s in %s mode...\n", cfg.Port, cfg.AppEnv)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

