package vehicle

import "time"

type FuelCard struct {
	ID              int64      `json:"id"`
	CardNo          string     `json:"card_no"`
	FuelCompanyID   int64      `json:"fuel_company_id"`
	FuelCompanyName string     `json:"fuel_company_name"`
	PinNumber       string     `json:"pin_number"`
	VehicleID       *int64     `json:"vehicle_id"`
	PlateNumber     *string    `json:"plate_number"`
	Status          bool       `json:"status"`
	CreatedBy       string     `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedBy       *string    `json:"updated_by"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
