package vehicle

type CondoPostalCode struct {
	ID         int64  `json:"id"`
	CondoID    int64  `json:"condo_id"`
	CondoName  string `json:"condo_name"`
	PostalCode string `json:"postal_code"`
}
