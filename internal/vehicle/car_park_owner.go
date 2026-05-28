package vehicle

import (
	"time"
)

type CarParkOwner struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Status    bool       `json:"status"`
	CreatedBy string     `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedBy *string    `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}
