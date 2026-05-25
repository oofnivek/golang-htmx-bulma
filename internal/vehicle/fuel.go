package vehicle

import "time"

type VehicleFuel struct {
	ID               int64      `json:"id"`
	Status           bool       `json:"status"`
	VehicleMakeID    int64      `json:"vehicle_make_id"`
	VehicleModelID   int64      `json:"vehicle_model_id"`
	FuelTypeID       int64      `json:"fuel_type_id"`
	VehicleMakeName  string     `json:"vehicle_make_name"`
	VehicleModelName string     `json:"vehicle_model_name"`
	FuelTypeName     string     `json:"fuel_type_name"`
	FuelTankSize     float64    `json:"fuel_tank_size"`
	FuelConsumption  float64    `json:"fuel_consumption"`
	CreatedBy        string     `json:"created_by"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedBy        *string    `json:"updated_by"`
	UpdatedAt        *time.Time `json:"updated_at"`
}
