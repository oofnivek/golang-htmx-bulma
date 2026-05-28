package vehicle

import "time"

type CarParkLot struct {
	ID          int64      `json:"id"`
	CarParkID   int64      `json:"car_park_id"`
	CarParkName string     `json:"car_park_name"`
	LotNumber   string     `json:"lot_number"`
	Level       string     `json:"level"`
	Status      bool       `json:"status"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedBy   *string    `json:"updated_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
