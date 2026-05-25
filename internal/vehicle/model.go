package vehicle

import "time"

type VehicleModel struct {
	ID              int64      `json:"id"`
	VehicleTypeID   int64      `json:"vehicle_type_id"`
	VehicleMakeID   int64      `json:"vehicle_make_id"`
	VehicleTypeName string     `json:"vehicle_type_name"`
	VehicleMakeName string     `json:"vehicle_make_name"`
	Name            string     `json:"name"`
	Status          bool       `json:"status"`
	CreatedBy       string     `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedBy       *string    `json:"updated_by"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
