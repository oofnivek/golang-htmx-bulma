package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"golang-htmx-bulma/internal/config"
	"golang-htmx-bulma/internal/db"
	"golang-htmx-bulma/internal/http/handlers/api"
	"golang-htmx-bulma/internal/http/handlers/web"
	"golang-htmx-bulma/internal/http/routes"
	"golang-htmx-bulma/internal/telemetry"
	"golang-htmx-bulma/internal/user"
	"golang-htmx-bulma/internal/vehicle"
	"golang-htmx-bulma/internal/view"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.Load()

	// Read application role: monolith (default), user-service, vehicle-service, web-view
	appRole := os.Getenv("APP_ROLE")
	if appRole == "" {
		appRole = "monolith"
	}

	// Initialize OpenTelemetry Logging
	ctx := context.Background()
	logProvider, err := telemetry.InitLogger(ctx, "fleet-management-system-"+appRole, cfg.AppEnv)
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry logging: %v\n", err)
	}
	defer func() {
		if err := logProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down OpenTelemetry logger provider: %v\n", err)
		}
	}()

	// Initialize Gin
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	if cfg.JWTSigningKey == "" {
		log.Fatal("JWT_SIGNING_KEY environment variable is not set")
	}

	// Wires variables
	var (
		vcHandler   *web.VehicleColorHandler
		vmHandler   *web.VehicleMakeHandler
		vtHandler   *web.VehicleTypeHandler
		vsHandler   *web.VehicleStatusHandler
		roleHandler *web.RoleHandler
		userHandler *web.UserHandler
		authHandler *web.AuthHandler

		userAPI *api.UserAPIHandler
		roleAPI *api.RoleAPIHandler
		authAPI *api.AuthAPIHandler
		vcAPI   *api.VehicleColorAPIHandler
		vmAPI   *api.VehicleMakeAPIHandler
		vsAPI   *api.VehicleStatusAPIHandler

		port string = cfg.Port
	)

	switch appRole {
	case "user-service":
		log.Println("Booting in USER-SERVICE role...")
		if os.Getenv("PORT") != "" {
			port = os.Getenv("PORT")
		} else {
			port = "8081" // Default user service port
		}

		// Initialize User DB only
		userDatabase, err := db.InitDB(cfg.UserDBDSN)
		if err != nil {
			log.Fatalf("Failed to connect to User MySQL: %v\n", err)
		}
		defer userDatabase.Close()

		// Wire services
		userRepo := user.NewUserRepository(userDatabase)
		userSvc := user.NewLocalUserService(userRepo)

		roleRepo := user.NewRoleRepository(userDatabase)
		roleSvc := user.NewLocalRoleService(roleRepo)

		authSvc := user.NewLocalAuthService(userSvc, cfg.JWTSigningKey)

		// Set up API handlers
		userAPI = api.NewUserAPIHandler(userSvc)
		roleAPI = api.NewRoleAPIHandler(roleSvc)
		authAPI = api.NewAuthAPIHandler(authSvc)

	case "vehicle-service":
		log.Println("Booting in VEHICLE-SERVICE role...")
		if os.Getenv("PORT") != "" {
			port = os.Getenv("PORT")
		} else {
			port = "8082" // Default vehicle service port
		}

		// Initialize Vehicle DB only
		vehicleDatabase, err := db.InitDB(cfg.VehicleDBDSN)
		if err != nil {
			log.Fatalf("Failed to connect to Vehicle MySQL: %v\n", err)
		}
		defer vehicleDatabase.Close()

		// Wire services
		vcRepo := vehicle.NewVehicleColorRepository(vehicleDatabase)
		vcSvc := vehicle.NewVehicleColorService(vcRepo)

		vmRepo := vehicle.NewVehicleMakeRepository(vehicleDatabase)
		vmSvc := vehicle.NewVehicleMakeService(vmRepo)

		vsRepo := vehicle.NewVehicleStatusRepository(vehicleDatabase)
		vsSvc := vehicle.NewVehicleStatusService(vsRepo)

		// Set up API handlers
		vcAPI = api.NewVehicleColorAPIHandler(vcSvc)
		vmAPI = api.NewVehicleMakeAPIHandler(vmSvc)
		vsAPI = api.NewVehicleStatusAPIHandler(vsSvc)

	case "web-view":
		log.Println("Booting in WEB-VIEW role...")
		if os.Getenv("PORT") != "" {
			port = os.Getenv("PORT")
		} else {
			port = "8080" // Default web-view port
		}

		// Wires remote services
		var (
			userSvc user.UserService
			roleSvc user.RoleService
			authSvc user.AuthService
		)

		userServiceURL := os.Getenv("USER_SERVICE_URL")
		if userServiceURL != "" {
			log.Printf("Using remote User Service at: %s\n", userServiceURL)
			userSvc = user.NewRemoteUserService(userServiceURL)
			roleSvc = user.NewRemoteRoleService(userServiceURL)
			authSvc = user.NewRemoteAuthService(userServiceURL)
		} else {
			log.Println("WARNING: USER_SERVICE_URL is not set. Falling back to local User DB access.")
			userDatabase, err := db.InitDB(cfg.UserDBDSN)
			if err != nil {
				log.Fatalf("Failed to connect to fallback User MySQL: %v\n", err)
			}
			defer userDatabase.Close()

			userRepo := user.NewUserRepository(userDatabase)
			userSvc = user.NewLocalUserService(userRepo)

			roleRepo := user.NewRoleRepository(userDatabase)
			roleSvc = user.NewLocalRoleService(roleRepo)

			authSvc = user.NewLocalAuthService(userSvc, cfg.JWTSigningKey)
		}

		// Wires Vehicle Services (remote or fallback local)
		var (
			vcSvc vehicle.VehicleColorService
			vmSvc vehicle.VehicleMakeService
			vtSvc vehicle.VehicleTypeService
		)

		var vsSvc vehicle.VehicleStatusService

		vehicleServiceURL := os.Getenv("VEHICLE_SERVICE_URL")
		if vehicleServiceURL != "" {
			log.Printf("Using remote Vehicle Service at: %s\n", vehicleServiceURL)
			vcSvc = vehicle.NewRemoteVehicleColorService(vehicleServiceURL)
			vmSvc = vehicle.NewRemoteVehicleMakeService(vehicleServiceURL)
			vtSvc = vehicle.NewRemoteVehicleTypeService(vehicleServiceURL)
			vsSvc = vehicle.NewRemoteVehicleStatusService(vehicleServiceURL)
		} else {
			log.Println("WARNING: VEHICLE_SERVICE_URL is not set. Falling back to local Vehicle DB access.")
			vehicleDatabase, err := db.InitDB(cfg.VehicleDBDSN)
			if err != nil {
				log.Fatalf("Failed to connect to fallback Vehicle MySQL: %v\n", err)
			}
			defer vehicleDatabase.Close()

			vcRepo := vehicle.NewVehicleColorRepository(vehicleDatabase)
			vcSvc = vehicle.NewVehicleColorService(vcRepo)

			vmRepo := vehicle.NewVehicleMakeRepository(vehicleDatabase)
			vmSvc = vehicle.NewVehicleMakeService(vmRepo)

			vtRepo := vehicle.NewVehicleTypeRepository(vehicleDatabase)
			vtSvc = vehicle.NewVehicleTypeService(vtRepo)

			vsRepo := vehicle.NewVehicleStatusRepository(vehicleDatabase)
			vsSvc = vehicle.NewVehicleStatusService(vsRepo)
		}

		// Set up Web Handlers
		vcHandler = web.NewVehicleColorHandler(vcSvc)
		vmHandler = web.NewVehicleMakeHandler(vmSvc)
		vtHandler = web.NewVehicleTypeHandler(vtSvc)
		vsHandler = web.NewVehicleStatusHandler(vsSvc)
		roleHandler = web.NewRoleHandler(roleSvc)
		userHandler = web.NewUserHandler(userSvc, roleSvc)
		authHandler = web.NewAuthHandler(authSvc)

		// Set up templates
		r.HTMLRender = view.NewRenderer("templates")

	case "monolith":
		fallthrough
	default:
		log.Println("Booting in MONOLITH role...")

		// Initialize both databases
		vehicleDatabase, err := db.InitDB(cfg.VehicleDBDSN)
		if err != nil {
			log.Fatalf("Failed to connect to Vehicle MySQL: %v\n", err)
		}
		defer vehicleDatabase.Close()

		userDatabase, err := db.InitDB(cfg.UserDBDSN)
		if err != nil {
			log.Fatalf("Failed to connect to User MySQL: %v\n", err)
		}
		defer userDatabase.Close()

		// Wire Vehicle services/handlers
		vcRepo := vehicle.NewVehicleColorRepository(vehicleDatabase)
		vcSvc := vehicle.NewVehicleColorService(vcRepo)
		vcHandler = web.NewVehicleColorHandler(vcSvc)

		vmRepo := vehicle.NewVehicleMakeRepository(vehicleDatabase)
		vmSvc := vehicle.NewVehicleMakeService(vmRepo)
		vmHandler = web.NewVehicleMakeHandler(vmSvc)

		vtRepo := vehicle.NewVehicleTypeRepository(vehicleDatabase)
		vtSvc := vehicle.NewVehicleTypeService(vtRepo)
		vtHandler = web.NewVehicleTypeHandler(vtSvc)

		vsRepo := vehicle.NewVehicleStatusRepository(vehicleDatabase)
		vsSvc := vehicle.NewVehicleStatusService(vsRepo)
		vsHandler = web.NewVehicleStatusHandler(vsSvc)

		// Wire User services/handlers (local implementations)
		userRepo := user.NewUserRepository(userDatabase)
		userSvc := user.NewLocalUserService(userRepo)

		roleRepo := user.NewRoleRepository(userDatabase)
		roleSvc := user.NewLocalRoleService(roleRepo)

		authSvc := user.NewLocalAuthService(userSvc, cfg.JWTSigningKey)

		// Set up Web Handlers
		roleHandler = web.NewRoleHandler(roleSvc)
		userHandler = web.NewUserHandler(userSvc, roleSvc)
		authHandler = web.NewAuthHandler(authSvc)

		// Set up API Handlers
		userAPI = api.NewUserAPIHandler(userSvc)
		roleAPI = api.NewRoleAPIHandler(roleSvc)
		authAPI = api.NewAuthAPIHandler(authSvc)
		vcAPI = api.NewVehicleColorAPIHandler(vcSvc)
		vmAPI = api.NewVehicleMakeAPIHandler(vmSvc)
		vsAPI = api.NewVehicleStatusAPIHandler(vsSvc)

		// Set up templates
		r.HTMLRender = view.NewRenderer("templates")
	}

	// Register routes
	routes.RegisterRoutes(
		r,
		vcHandler,
		vmHandler,
		vtHandler,
		vsHandler,
		roleHandler,
		userHandler,
		authHandler,
		userAPI,
		roleAPI,
		authAPI,
		vcAPI,
		vmAPI,
		vsAPI,
		cfg.JWTSigningKey,
	)

	// Start server
	log.Printf("Server starting on port %s in %s mode...\n", port, cfg.AppEnv)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
