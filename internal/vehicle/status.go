package vehicle

type VehicleStatus struct {
	ID        int64  `json:"id"`
	IsActive  bool   `json:"is_active"`
	Substatus string `json:"substatus"`
}
