package vehicle

import "time"

type VehicleGlobalSetting struct {
	ID          int64      `json:"id"`
	Key         string     `json:"key"`
	Value       string     `json:"value"`
	Remark      *string    `json:"remark"`
	CountryCode *string    `json:"country_code"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedBy   *string    `json:"updated_by"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
