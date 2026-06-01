package vehicle

type RegionalInfo struct {
	PostalCode string `json:"postal_code"`
	Region     string `json:"region"`
	EstateID   int64  `json:"estate_id"`
	EstateName string `json:"estate_name"`
}
