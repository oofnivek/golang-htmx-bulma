package vehicle

import "time"

type VehicleService interface {
	ListAll() ([]Vehicle, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]Vehicle, int, error)
	FindByID(id int64) (*Vehicle, error)
	CreateVehicle(
		vehicleMakeID, vehicleModelID, vehicleTypeID, fuelTypeID, vehicleColorID int64,
		description, plateNumber, iuNumber, chassisNumber, engineNumber, bootSpace *string,
		numSeats int,
		carParkID, assetOwnerID, vehicleStatusID int64,
		lastServiceDate, lastCleanedDate, activeFrom, activeTo *time.Time,
		lastServiceMileage, currentMileage, currentFuelLevel *int,
		user string,
	) (*Vehicle, error)
	UpdateVehicle(
		id int64,
		vehicleMakeID, vehicleModelID, vehicleTypeID, fuelTypeID, vehicleColorID int64,
		description, plateNumber, iuNumber, chassisNumber, engineNumber, bootSpace *string,
		numSeats int,
		carParkID, assetOwnerID, vehicleStatusID int64,
		lastServiceDate, lastCleanedDate, activeFrom, activeTo *time.Time,
		lastServiceMileage, currentMileage, currentFuelLevel *int,
		user string,
	) (*Vehicle, error)
	DeleteVehicle(id int64) error
}

type vehicleService struct {
	repo VehicleRepository
}

func NewVehicleService(repo VehicleRepository) VehicleService {
	return &vehicleService{repo: repo}
}

func (s *vehicleService) ListAll() ([]Vehicle, error) {
	return s.repo.GetAll()
}

func (s *vehicleService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]Vehicle, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	vehicles, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return vehicles, total, nil
}

func (s *vehicleService) FindByID(id int64) (*Vehicle, error) {
	return s.repo.GetByID(id)
}

func (s *vehicleService) CreateVehicle(
	vehicleMakeID, vehicleModelID, vehicleTypeID, fuelTypeID, vehicleColorID int64,
	description, plateNumber, iuNumber, chassisNumber, engineNumber, bootSpace *string,
	numSeats int,
	carParkID, assetOwnerID, vehicleStatusID int64,
	lastServiceDate, lastCleanedDate, activeFrom, activeTo *time.Time,
	lastServiceMileage, currentMileage, currentFuelLevel *int,
	user string,
) (*Vehicle, error) {
	now := time.Now().UTC()
	v := &Vehicle{
		VehicleMakeID:      vehicleMakeID,
		VehicleModelID:     vehicleModelID,
		VehicleTypeID:      vehicleTypeID,
		FuelTypeID:         fuelTypeID,
		VehicleColorID:     vehicleColorID,
		Description:        description,
		PlateNumber:        plateNumber,
		IUNumber:           iuNumber,
		ChassisNumber:      chassisNumber,
		EngineNumber:       engineNumber,
		NumSeats:           numSeats,
		BootSpace:          bootSpace,
		CarParkID:          carParkID,
		AssetOwnerID:       assetOwnerID,
		VehicleStatusID:    vehicleStatusID,
		LastServiceDate:    lastServiceDate,
		LastCleanedDate:    lastCleanedDate,
		LastServiceMileage: lastServiceMileage,
		CurrentMileage:     currentMileage,
		CurrentFuelLevel:   currentFuelLevel,
		ActiveFrom:         activeFrom,
		ActiveTo:           activeTo,
		CreatedBy:          user,
		CreatedAt:          now,
		UpdatedBy:          &user,
		UpdatedAt:          &now,
	}
	if err := s.repo.Create(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *vehicleService) UpdateVehicle(
	id int64,
	vehicleMakeID, vehicleModelID, vehicleTypeID, fuelTypeID, vehicleColorID int64,
	description, plateNumber, iuNumber, chassisNumber, engineNumber, bootSpace *string,
	numSeats int,
	carParkID, assetOwnerID, vehicleStatusID int64,
	lastServiceDate, lastCleanedDate, activeFrom, activeTo *time.Time,
	lastServiceMileage, currentMileage, currentFuelLevel *int,
	user string,
) (*Vehicle, error) {
	v, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}

	now := time.Now().UTC()
	v.VehicleMakeID = vehicleMakeID
	v.VehicleModelID = vehicleModelID
	v.VehicleTypeID = vehicleTypeID
	v.FuelTypeID = fuelTypeID
	v.VehicleColorID = vehicleColorID
	v.Description = description
	v.PlateNumber = plateNumber
	v.IUNumber = iuNumber
	v.ChassisNumber = chassisNumber
	v.EngineNumber = engineNumber
	v.NumSeats = numSeats
	v.BootSpace = bootSpace
	v.CarParkID = carParkID
	v.AssetOwnerID = assetOwnerID
	v.VehicleStatusID = vehicleStatusID
	v.LastServiceDate = lastServiceDate
	v.LastCleanedDate = lastCleanedDate
	v.LastServiceMileage = lastServiceMileage
	v.CurrentMileage = currentMileage
	v.CurrentFuelLevel = currentFuelLevel
	v.ActiveFrom = activeFrom
	v.ActiveTo = activeTo
	v.UpdatedBy = &user
	v.UpdatedAt = &now

	if err := s.repo.Update(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *vehicleService) DeleteVehicle(id int64) error {
	return s.repo.Delete(id)
}
