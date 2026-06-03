package vehicle

type CondoCarPark struct {
	ID          int64  `json:"id"`
	CondoID     int64  `json:"condo_id"`
	CondoName   string `json:"condo_name"`
	CarParkID   int64  `json:"car_park_id"`
	CarParkName string `json:"car_park_name"`
}
