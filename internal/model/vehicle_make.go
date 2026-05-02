package model

import (
	"time"
)

type VehicleMake struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Status    bool       `json:"status"` // tinyint(1)
	CreatedBy string     `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedBy *string    `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}
