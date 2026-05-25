package vehicle

import "time"

type VehicleFuelService interface {
	ListAll() ([]VehicleFuel, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleFuel, int, error)
	FindByID(id int64) (*VehicleFuel, error)
	CreateFuel(vehicleMakeID, vehicleModelID, fuelTypeID int64, fuelTankSize, fuelConsumption float64, status bool, user string) (*VehicleFuel, error)
	UpdateFuel(id, vehicleMakeID, vehicleModelID, fuelTypeID int64, fuelTankSize, fuelConsumption float64, status bool, user string) (*VehicleFuel, error)
	DeleteFuel(id int64) error
}

type vehicleFuelService struct {
	repo VehicleFuelRepository
}

func NewVehicleFuelService(repo VehicleFuelRepository) VehicleFuelService {
	return &vehicleFuelService{repo: repo}
}

func (s *vehicleFuelService) ListAll() ([]VehicleFuel, error) {
	return s.repo.GetAll()
}

func (s *vehicleFuelService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleFuel, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	fuels, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}
	return fuels, total, nil
}

func (s *vehicleFuelService) FindByID(id int64) (*VehicleFuel, error) {
	return s.repo.GetByID(id)
}

func (s *vehicleFuelService) CreateFuel(vehicleMakeID, vehicleModelID, fuelTypeID int64, fuelTankSize, fuelConsumption float64, status bool, user string) (*VehicleFuel, error) {
	now := time.Now().UTC()
	f := &VehicleFuel{
		Status:          status,
		VehicleMakeID:   vehicleMakeID,
		VehicleModelID:  vehicleModelID,
		FuelTypeID:      fuelTypeID,
		FuelTankSize:    fuelTankSize,
		FuelConsumption: fuelConsumption,
		CreatedBy:       user,
		CreatedAt:       now,
		UpdatedBy:       &user,
		UpdatedAt:       &now,
	}
	if err := s.repo.Create(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *vehicleFuelService) UpdateFuel(id, vehicleMakeID, vehicleModelID, fuelTypeID int64, fuelTankSize, fuelConsumption float64, status bool, user string) (*VehicleFuel, error) {
	f, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, nil
	}

	now := time.Now().UTC()
	f.Status = status
	f.VehicleMakeID = vehicleMakeID
	f.VehicleModelID = vehicleModelID
	f.FuelTypeID = fuelTypeID
	f.FuelTankSize = fuelTankSize
	f.FuelConsumption = fuelConsumption
	f.UpdatedBy = &user
	f.UpdatedAt = &now

	if err := s.repo.Update(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *vehicleFuelService) DeleteFuel(id int64) error {
	return s.repo.Delete(id)
}
