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
		ftHandler   *web.FuelTypeHandler
		vmdlHandler *web.VehicleModelHandler
		vfHandler   *web.VehicleFuelHandler
		fcHandler   *web.FuelCompanyHandler
		caoHandler  *web.CarAssetOwnerHandler
		cpoHandler  *web.CarParkOwnerHandler
		cpHandler   *web.CarParkHandler
		cplHandler  *web.CarParkLotHandler
		roleHandler *web.RoleHandler
		userHandler *web.UserHandler
		authHandler *web.AuthHandler

		userAPI  *api.UserAPIHandler
		roleAPI  *api.RoleAPIHandler
		authAPI  *api.AuthAPIHandler
		vcAPI    *api.VehicleColorAPIHandler
		vmAPI    *api.VehicleMakeAPIHandler
		vsAPI    *api.VehicleStatusAPIHandler
		ftAPI    *api.FuelTypeAPIHandler
		vmdlAPI  *api.VehicleModelAPIHandler
		vfAPI    *api.VehicleFuelAPIHandler
		fcAPI    *api.FuelCompanyAPIHandler
		caoAPI  *api.CarAssetOwnerAPIHandler
		cpoAPI  *api.CarParkOwnerAPIHandler
		cpAPI   *api.CarParkAPIHandler
		cplAPI  *api.CarParkLotAPIHandler
		estateHandler *web.EstateHandler
		estateAPI     *api.EstateAPIHandler
		riHandler     *web.RegionalInfoHandler
		riAPI         *api.RegionalInfoAPIHandler
		vHandler      *web.VehicleHandler
		vAPI          *api.VehicleAPIHandler
		fuelCardHandler *web.FuelCardHandler
		fuelCardAPI     *api.FuelCardAPIHandler
		condoHandler        *web.CondoHandler
		condoAPI            *api.CondoAPIHandler
		condoCarParkHandler *web.CondoCarParkHandler
		condoCarParkAPI     *api.CondoCarParkAPIHandler
		condoPostalCodeHandler *web.CondoPostalCodeHandler
		condoPostalCodeAPI     *api.CondoPostalCodeAPIHandler

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

		vtRepo := vehicle.NewVehicleTypeRepository(vehicleDatabase)
		vtSvc := vehicle.NewVehicleTypeService(vtRepo)

		vsRepo := vehicle.NewVehicleStatusRepository(vehicleDatabase)
		vsSvc := vehicle.NewVehicleStatusService(vsRepo)

		ftRepo := vehicle.NewFuelTypeRepository(vehicleDatabase)
		ftSvc := vehicle.NewFuelTypeService(ftRepo)

		vmdlRepo := vehicle.NewVehicleModelRepository(vehicleDatabase)
		vmdlSvc := vehicle.NewVehicleModelService(vmdlRepo)

		vfRepo := vehicle.NewVehicleFuelRepository(vehicleDatabase)
		vfSvc := vehicle.NewVehicleFuelService(vfRepo)

		fcRepo := vehicle.NewFuelCompanyRepository(vehicleDatabase)
		fcSvc := vehicle.NewFuelCompanyService(fcRepo)

		caoRepo := vehicle.NewCarAssetOwnerRepository(vehicleDatabase)
		caoSvc := vehicle.NewCarAssetOwnerService(caoRepo)

		cpoRepo := vehicle.NewCarParkOwnerRepository(vehicleDatabase)
		cpoSvc := vehicle.NewCarParkOwnerService(cpoRepo)

		cpRepo := vehicle.NewCarParkRepository(vehicleDatabase)
		cpSvc := vehicle.NewCarParkService(cpRepo)

		cplRepo := vehicle.NewCarParkLotRepository(vehicleDatabase)
		cplSvc := vehicle.NewCarParkLotService(cplRepo)

		estateRepo := vehicle.NewEstateRepository(vehicleDatabase)
		estateSvc := vehicle.NewEstateService(estateRepo)

		riRepo := vehicle.NewRegionalInfoRepository(vehicleDatabase)
		riSvc := vehicle.NewRegionalInfoService(riRepo)

		vRepo := vehicle.NewVehicleRepository(vehicleDatabase)
		vSvc := vehicle.NewVehicleService(vRepo)

		// Set up API handlers
		vcAPI = api.NewVehicleColorAPIHandler(vcSvc)
		vmAPI = api.NewVehicleMakeAPIHandler(vmSvc)
		vsAPI = api.NewVehicleStatusAPIHandler(vsSvc)
		ftAPI = api.NewFuelTypeAPIHandler(ftSvc)
		vmdlAPI = api.NewVehicleModelAPIHandler(vmdlSvc)
		vfAPI = api.NewVehicleFuelAPIHandler(vfSvc)
		fcAPI = api.NewFuelCompanyAPIHandler(fcSvc)
		caoAPI = api.NewCarAssetOwnerAPIHandler(caoSvc)
		cpoAPI = api.NewCarParkOwnerAPIHandler(cpoSvc)
		cpAPI = api.NewCarParkAPIHandler(cpSvc)
		cplAPI = api.NewCarParkLotAPIHandler(cplSvc)
		estateAPI = api.NewEstateAPIHandler(estateSvc)
		riAPI = api.NewRegionalInfoAPIHandler(riSvc)
		vAPI = api.NewVehicleAPIHandler(vSvc)

		fuelCardRepo := vehicle.NewFuelCardRepository(vehicleDatabase)
		fuelCardSvc := vehicle.NewFuelCardService(fuelCardRepo)
		fuelCardAPI = api.NewFuelCardAPIHandler(fuelCardSvc)

		condoRepo := vehicle.NewCondoRepository(vehicleDatabase)
		condoSvc := vehicle.NewCondoService(condoRepo)
		condoAPI = api.NewCondoAPIHandler(condoSvc)

		condoCarParkRepo := vehicle.NewCondoCarParkRepository(vehicleDatabase)
		condoCarParkSvc := vehicle.NewCondoCarParkService(condoCarParkRepo)
		condoCarParkAPI = api.NewCondoCarParkAPIHandler(condoCarParkSvc)

		condoPostalCodeRepo := vehicle.NewCondoPostalCodeRepository(vehicleDatabase)
		condoPostalCodeSvc := vehicle.NewCondoPostalCodeService(condoPostalCodeRepo)
		condoPostalCodeAPI = api.NewCondoPostalCodeAPIHandler(condoPostalCodeSvc)

		// Set up Web Handlers
		vHandler = web.NewVehicleHandler(vSvc, vmSvc, vmdlSvc, vtSvc, ftSvc, vcSvc, cpSvc, caoSvc, vsSvc)

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
			vcSvc   vehicle.VehicleColorService
			vmSvc   vehicle.VehicleMakeService
			vtSvc   vehicle.VehicleTypeService
			vmdlSvc vehicle.VehicleModelService
			vfSvc   vehicle.VehicleFuelService
		)

		var vsSvc vehicle.VehicleStatusService

		var ftSvc vehicle.FuelTypeService

		var fcSvc vehicle.FuelCompanyService

		var caoSvc vehicle.CarAssetOwnerService

		var cpoSvc vehicle.CarParkOwnerService

		var cpSvc vehicle.CarParkService

		var cplSvc vehicle.CarParkLotService

		var estateSvc vehicle.EstateService

		var riSvc vehicle.RegionalInfoService

		var fuelCardSvc vehicle.FuelCardService

		var condoSvc vehicle.CondoService

		var condoCarParkSvc vehicle.CondoCarParkService

		var condoPostalCodeSvc vehicle.CondoPostalCodeService

		vehicleServiceURL := os.Getenv("VEHICLE_SERVICE_URL")
		var vSvc vehicle.VehicleService

		if vehicleServiceURL != "" {
			log.Printf("Using remote Vehicle Service at: %s\n", vehicleServiceURL)
			vcSvc = vehicle.NewRemoteVehicleColorService(vehicleServiceURL)
			vmSvc = vehicle.NewRemoteVehicleMakeService(vehicleServiceURL)
			vtSvc = vehicle.NewRemoteVehicleTypeService(vehicleServiceURL)
			vsSvc = vehicle.NewRemoteVehicleStatusService(vehicleServiceURL)
			ftSvc = vehicle.NewRemoteFuelTypeService(vehicleServiceURL)
			vmdlSvc = vehicle.NewRemoteVehicleModelService(vehicleServiceURL)
			vfSvc = vehicle.NewRemoteVehicleFuelService(vehicleServiceURL)
			fcSvc = vehicle.NewRemoteFuelCompanyService(vehicleServiceURL)
			caoSvc = vehicle.NewRemoteCarAssetOwnerService(vehicleServiceURL)
			cpoSvc = vehicle.NewRemoteCarParkOwnerService(vehicleServiceURL)
			cpSvc = vehicle.NewRemoteCarParkService(vehicleServiceURL)
			cplSvc = vehicle.NewRemoteCarParkLotService(vehicleServiceURL)
			estateSvc = vehicle.NewRemoteEstateService(vehicleServiceURL)
			riSvc = vehicle.NewRemoteRegionalInfoService(vehicleServiceURL)
			vSvc = vehicle.NewRemoteVehicleService(vehicleServiceURL)
			fuelCardSvc = vehicle.NewRemoteFuelCardService(vehicleServiceURL)
			condoSvc = vehicle.NewRemoteCondoService(vehicleServiceURL)
			condoCarParkSvc = vehicle.NewRemoteCondoCarParkService(vehicleServiceURL)
			condoPostalCodeSvc = vehicle.NewRemoteCondoPostalCodeService(vehicleServiceURL)
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

			ftRepo := vehicle.NewFuelTypeRepository(vehicleDatabase)
			ftSvc = vehicle.NewFuelTypeService(ftRepo)

			vmdlRepo := vehicle.NewVehicleModelRepository(vehicleDatabase)
			vmdlSvc = vehicle.NewVehicleModelService(vmdlRepo)

			vfRepo := vehicle.NewVehicleFuelRepository(vehicleDatabase)
			vfSvc = vehicle.NewVehicleFuelService(vfRepo)

			fcRepo := vehicle.NewFuelCompanyRepository(vehicleDatabase)
			fcSvc = vehicle.NewFuelCompanyService(fcRepo)

			caoRepo := vehicle.NewCarAssetOwnerRepository(vehicleDatabase)
			caoSvc = vehicle.NewCarAssetOwnerService(caoRepo)

			cpoRepo := vehicle.NewCarParkOwnerRepository(vehicleDatabase)
			cpoSvc = vehicle.NewCarParkOwnerService(cpoRepo)

			cpRepo := vehicle.NewCarParkRepository(vehicleDatabase)
			cpSvc = vehicle.NewCarParkService(cpRepo)

			cplRepo := vehicle.NewCarParkLotRepository(vehicleDatabase)
			cplSvc = vehicle.NewCarParkLotService(cplRepo)

			estateRepo := vehicle.NewEstateRepository(vehicleDatabase)
			estateSvc = vehicle.NewEstateService(estateRepo)

			riRepo := vehicle.NewRegionalInfoRepository(vehicleDatabase)
			riSvc = vehicle.NewRegionalInfoService(riRepo)

			vRepo := vehicle.NewVehicleRepository(vehicleDatabase)
			vSvc = vehicle.NewVehicleService(vRepo)

			fuelCardRepo := vehicle.NewFuelCardRepository(vehicleDatabase)
			fuelCardSvc = vehicle.NewFuelCardService(fuelCardRepo)

			condoRepo := vehicle.NewCondoRepository(vehicleDatabase)
			condoSvc = vehicle.NewCondoService(condoRepo)

			condoCarParkRepo := vehicle.NewCondoCarParkRepository(vehicleDatabase)
			condoCarParkSvc = vehicle.NewCondoCarParkService(condoCarParkRepo)

			condoPostalCodeRepo := vehicle.NewCondoPostalCodeRepository(vehicleDatabase)
			condoPostalCodeSvc = vehicle.NewCondoPostalCodeService(condoPostalCodeRepo)
		}

		// Set up Web Handlers
		vcHandler = web.NewVehicleColorHandler(vcSvc)
		vmHandler = web.NewVehicleMakeHandler(vmSvc)
		vtHandler = web.NewVehicleTypeHandler(vtSvc)
		vsHandler = web.NewVehicleStatusHandler(vsSvc)
		ftHandler = web.NewFuelTypeHandler(ftSvc)
		vmdlHandler = web.NewVehicleModelHandler(vmdlSvc, vtSvc, vmSvc)
		vfHandler = web.NewVehicleFuelHandler(vfSvc, vmSvc, vmdlSvc, ftSvc)
		fcHandler = web.NewFuelCompanyHandler(fcSvc)
		caoHandler = web.NewCarAssetOwnerHandler(caoSvc)
		cpoHandler = web.NewCarParkOwnerHandler(cpoSvc)
		cpHandler = web.NewCarParkHandler(cpSvc, cpoSvc)
		cplHandler = web.NewCarParkLotHandler(cplSvc, cpSvc)
		estateHandler = web.NewEstateHandler(estateSvc)
		riHandler = web.NewRegionalInfoHandler(riSvc, estateSvc)
		vHandler = web.NewVehicleHandler(vSvc, vmSvc, vmdlSvc, vtSvc, ftSvc, vcSvc, cpSvc, caoSvc, vsSvc)
		fuelCardHandler = web.NewFuelCardHandler(fuelCardSvc, fcSvc, vSvc)
		condoHandler = web.NewCondoHandler(condoSvc)
		condoCarParkHandler = web.NewCondoCarParkHandler(condoCarParkSvc, condoSvc, cpSvc)
		condoPostalCodeHandler = web.NewCondoPostalCodeHandler(condoPostalCodeSvc, condoSvc)
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

		ftRepo := vehicle.NewFuelTypeRepository(vehicleDatabase)
		ftSvc := vehicle.NewFuelTypeService(ftRepo)
		ftHandler = web.NewFuelTypeHandler(ftSvc)

		vmdlRepo := vehicle.NewVehicleModelRepository(vehicleDatabase)
		vmdlSvc := vehicle.NewVehicleModelService(vmdlRepo)
		vmdlHandler = web.NewVehicleModelHandler(vmdlSvc, vtSvc, vmSvc)

		vfRepo := vehicle.NewVehicleFuelRepository(vehicleDatabase)
		vfSvc := vehicle.NewVehicleFuelService(vfRepo)
		vfHandler = web.NewVehicleFuelHandler(vfSvc, vmSvc, vmdlSvc, ftSvc)

		fcRepo := vehicle.NewFuelCompanyRepository(vehicleDatabase)
		fcSvc := vehicle.NewFuelCompanyService(fcRepo)
		fcHandler = web.NewFuelCompanyHandler(fcSvc)

		caoRepo := vehicle.NewCarAssetOwnerRepository(vehicleDatabase)
		caoSvc := vehicle.NewCarAssetOwnerService(caoRepo)
		caoHandler = web.NewCarAssetOwnerHandler(caoSvc)

		cpoRepo := vehicle.NewCarParkOwnerRepository(vehicleDatabase)
		cpoSvc := vehicle.NewCarParkOwnerService(cpoRepo)
		cpoHandler = web.NewCarParkOwnerHandler(cpoSvc)

		cpRepo := vehicle.NewCarParkRepository(vehicleDatabase)
		cpSvc := vehicle.NewCarParkService(cpRepo)
		cpHandler = web.NewCarParkHandler(cpSvc, cpoSvc)

		cplRepo := vehicle.NewCarParkLotRepository(vehicleDatabase)
		cplSvc := vehicle.NewCarParkLotService(cplRepo)
		cplHandler = web.NewCarParkLotHandler(cplSvc, cpSvc)

		estateRepo := vehicle.NewEstateRepository(vehicleDatabase)
		estateSvc := vehicle.NewEstateService(estateRepo)
		estateHandler = web.NewEstateHandler(estateSvc)

		riRepo := vehicle.NewRegionalInfoRepository(vehicleDatabase)
		riSvc := vehicle.NewRegionalInfoService(riRepo)
		riHandler = web.NewRegionalInfoHandler(riSvc, estateSvc)

		vRepo := vehicle.NewVehicleRepository(vehicleDatabase)
		vSvc := vehicle.NewVehicleService(vRepo)
		vHandler = web.NewVehicleHandler(vSvc, vmSvc, vmdlSvc, vtSvc, ftSvc, vcSvc, cpSvc, caoSvc, vsSvc)
		vAPI = api.NewVehicleAPIHandler(vSvc)

		fuelCardRepo := vehicle.NewFuelCardRepository(vehicleDatabase)
		fuelCardSvc := vehicle.NewFuelCardService(fuelCardRepo)
		fuelCardHandler = web.NewFuelCardHandler(fuelCardSvc, fcSvc, vSvc)
		fuelCardAPI = api.NewFuelCardAPIHandler(fuelCardSvc)

		condoRepo := vehicle.NewCondoRepository(vehicleDatabase)
		condoSvc := vehicle.NewCondoService(condoRepo)
		condoHandler = web.NewCondoHandler(condoSvc)
		condoAPI = api.NewCondoAPIHandler(condoSvc)

		condoCarParkRepo := vehicle.NewCondoCarParkRepository(vehicleDatabase)
		condoCarParkSvc := vehicle.NewCondoCarParkService(condoCarParkRepo)
		condoCarParkHandler = web.NewCondoCarParkHandler(condoCarParkSvc, condoSvc, cpSvc)
		condoCarParkAPI = api.NewCondoCarParkAPIHandler(condoCarParkSvc)

		condoPostalCodeRepo := vehicle.NewCondoPostalCodeRepository(vehicleDatabase)
		condoPostalCodeSvc := vehicle.NewCondoPostalCodeService(condoPostalCodeRepo)
		condoPostalCodeHandler = web.NewCondoPostalCodeHandler(condoPostalCodeSvc, condoSvc)
		condoPostalCodeAPI = api.NewCondoPostalCodeAPIHandler(condoPostalCodeSvc)

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
		ftAPI = api.NewFuelTypeAPIHandler(ftSvc)
		vmdlAPI = api.NewVehicleModelAPIHandler(vmdlSvc)
		vfAPI = api.NewVehicleFuelAPIHandler(vfSvc)
		fcAPI = api.NewFuelCompanyAPIHandler(fcSvc)
		caoAPI = api.NewCarAssetOwnerAPIHandler(caoSvc)
		cpoAPI = api.NewCarParkOwnerAPIHandler(cpoSvc)
		cpAPI = api.NewCarParkAPIHandler(cpSvc)
		cplAPI = api.NewCarParkLotAPIHandler(cplSvc)
		estateAPI = api.NewEstateAPIHandler(estateSvc)
		riAPI = api.NewRegionalInfoAPIHandler(riSvc)

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
		ftHandler,
		vmdlHandler,
		vfHandler,
		fcHandler,
		caoHandler,
		cpoHandler,
		cpHandler,
		cplHandler,
		estateHandler,
		riHandler,
		vHandler,
		fuelCardHandler,
		condoHandler,
		condoCarParkHandler,
		roleHandler,
		userHandler,
		authHandler,
		userAPI,
		roleAPI,
		authAPI,
		vcAPI,
		vmAPI,
		vsAPI,
		ftAPI,
		vmdlAPI,
		vfAPI,
		fcAPI,
		caoAPI,
		cpoAPI,
		cpAPI,
		cplAPI,
		estateAPI,
		riAPI,
		vAPI,
		fuelCardAPI,
		condoAPI,
		condoCarParkAPI,
		condoPostalCodeHandler,
		condoPostalCodeAPI,
		cfg.JWTSigningKey,
	)

	// Start server
	log.Printf("Server starting on port %s in %s mode...\n", port, cfg.AppEnv)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
