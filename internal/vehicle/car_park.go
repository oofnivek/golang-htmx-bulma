package vehicle

import "time"

type CarPark struct {
	ID               int64      `json:"id"`
	Name             string     `json:"name"`
	Description      *string    `json:"description"`
	PostalCode       string     `json:"postal_code"`
	Address          string     `json:"address"`
	Latitude         float64    `json:"latitude"`
	Longitude        float64    `json:"longitude"`
	CarParkOwnerID   int64      `json:"car_park_owner_id"`
	CarParkOwnerName string     `json:"car_park_owner_name"`
	ActiveFrom       *time.Time `json:"active_from"`
	ActiveTo         *time.Time `json:"active_to"`
	Status           bool       `json:"status"`
	CreatedBy        string     `json:"created_by"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedBy        *string    `json:"updated_by"`
	UpdatedAt        *time.Time `json:"updated_at"`
}
