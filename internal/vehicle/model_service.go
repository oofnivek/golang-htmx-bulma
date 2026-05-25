package vehicle

import "time"

type VehicleModelService interface {
	ListAll() ([]VehicleModel, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleModel, int, error)
	FindByID(id int64) (*VehicleModel, error)
	CreateModel(vehicleTypeID, vehicleMakeID int64, name string, status bool, user string) (*VehicleModel, error)
	UpdateModel(id, vehicleTypeID, vehicleMakeID int64, name string, status bool, user string) (*VehicleModel, error)
	DeleteModel(id int64) error
}

type vehicleModelService struct {
	repo VehicleModelRepository
}

func NewVehicleModelService(repo VehicleModelRepository) VehicleModelService {
	return &vehicleModelService{repo: repo}
}

func (s *vehicleModelService) ListAll() ([]VehicleModel, error) {
	return s.repo.GetAll()
}

func (s *vehicleModelService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleModel, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	models, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return models, total, nil
}

func (s *vehicleModelService) FindByID(id int64) (*VehicleModel, error) {
	return s.repo.GetByID(id)
}

func (s *vehicleModelService) CreateModel(vehicleTypeID, vehicleMakeID int64, name string, status bool, user string) (*VehicleModel, error) {
	now := time.Now()
	m := &VehicleModel{
		VehicleTypeID: vehicleTypeID,
		VehicleMakeID: vehicleMakeID,
		Name:          name,
		Status:        status,
		CreatedBy:     user,
		CreatedAt:     now,
		UpdatedBy:     &user,
		UpdatedAt:     &now,
	}
	if err := s.repo.Create(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *vehicleModelService) UpdateModel(id, vehicleTypeID, vehicleMakeID int64, name string, status bool, user string) (*VehicleModel, error) {
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	m.VehicleTypeID = vehicleTypeID
	m.VehicleMakeID = vehicleMakeID
	m.Name = name
	m.Status = status
	now := time.Now()
	m.UpdatedBy = &user
	m.UpdatedAt = &now

	if err := s.repo.Update(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *vehicleModelService) DeleteModel(id int64) error {
	return s.repo.Delete(id)
}
