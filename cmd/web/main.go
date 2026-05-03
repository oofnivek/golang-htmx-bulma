package main

import (
	"log"
	"net/http"

	"golang-htmx-bulma/internal/config"
	"golang-htmx-bulma/internal/db"
	"golang-htmx-bulma/internal/http/handlers/web"
	"golang-htmx-bulma/internal/http/routes"
	"golang-htmx-bulma/internal/repository"
	"golang-htmx-bulma/internal/service"
	"golang-htmx-bulma/internal/view"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.Load()

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
	vcRepo := repository.NewVehicleColorRepository(vehicleDatabase)
	vcSvc := service.NewVehicleColorService(vcRepo)
	vcHandler := web.NewVehicleColorHandler(vcSvc)

	vmRepo := repository.NewVehicleMakeRepository(vehicleDatabase)
	vmSvc := service.NewVehicleMakeService(vmRepo)
	vmHandler := web.NewVehicleMakeHandler(vmSvc)

	roleRepo := repository.NewRoleRepository(userDatabase)
	roleSvc := service.NewRoleService(roleRepo)
	roleHandler := web.NewRoleHandler(roleSvc)

	routes.RegisterRoutes(r, vcHandler, vmHandler, roleHandler)

	// Start server
	log.Printf("Server starting on port %s in %s mode...\n", cfg.Port, cfg.AppEnv)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

