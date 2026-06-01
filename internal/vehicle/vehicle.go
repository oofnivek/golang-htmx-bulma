package vehicle

import "time"

type Vehicle struct {
	ID                 int64      `json:"id"`
	VehicleMakeID      int64      `json:"vehicle_make_id"`
	VehicleMakeName    string     `json:"vehicle_make_name"`
	VehicleModelID     int64      `json:"vehicle_model_id"`
	VehicleModelName   string     `json:"vehicle_model_name"`
	VehicleTypeID      int64      `json:"vehicle_type_id"`
	VehicleTypeName    string     `json:"vehicle_type_name"`
	FuelTypeID         int64      `json:"fuel_type_id"`
	FuelTypeName       string     `json:"fuel_type_name"`
	VehicleColorID     int64      `json:"vehicle_color_id"`
	VehicleColorName   string     `json:"vehicle_color_name"`
	Description        *string    `json:"description"`
	PlateNumber        *string    `json:"plate_number"`
	IUNumber           *string    `json:"iu_number"`
	ChassisNumber      *string    `json:"chassis_number"`
	EngineNumber       *string    `json:"engine_number"`
	NumSeats           int        `json:"num_seats"`
	BootSpace          *string    `json:"boot_space"`
	CarParkID          int64      `json:"car_park_id"`
	CarParkName        string     `json:"car_park_name"`
	AssetOwnerID       int64      `json:"asset_owner_id"`
	AssetOwnerName     string     `json:"asset_owner_name"`
	LastServiceDate    *time.Time `json:"last_service_date"`
	LastCleanedDate    *time.Time `json:"last_cleaned_date"`
	LastServiceMileage *int       `json:"last_service_mileage"`
	CurrentMileage     *int       `json:"current_mileage"`
	CurrentFuelLevel   *int       `json:"current_fuel_level"`
	ActiveFrom         *time.Time `json:"active_from"`
	ActiveTo           *time.Time `json:"active_to"`
	VehicleStatusID    int64      `json:"vehicle_status_id"`
	VehicleStatusName  string     `json:"vehicle_status_name"`
	CreatedBy          string     `json:"created_by"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedBy          *string    `json:"updated_by"`
	UpdatedAt          *time.Time `json:"updated_at"`
}
