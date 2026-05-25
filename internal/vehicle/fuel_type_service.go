package vehicle

import "time"

type FuelTypeService interface {
	ListAll() ([]FuelType, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]FuelType, int, error)
	FindByID(id int64) (*FuelType, error)
	CreateFuelType(name string, status bool, user string) (*FuelType, error)
	UpdateFuelType(id int64, name string, status bool, user string) (*FuelType, error)
	DeleteFuelType(id int64) error
}

type fuelTypeService struct {
	repo FuelTypeRepository
}

func NewFuelTypeService(repo FuelTypeRepository) FuelTypeService {
	return &fuelTypeService{repo: repo}
}

func (s *fuelTypeService) ListAll() ([]FuelType, error) {
	return s.repo.GetAll()
}

func (s *fuelTypeService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]FuelType, int, error) {
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

func (s *fuelTypeService) FindByID(id int64) (*FuelType, error) {
	return s.repo.GetByID(id)
}

func (s *fuelTypeService) CreateFuelType(name string, status bool, user string) (*FuelType, error) {
	now := time.Now().UTC()
	f := &FuelType{
		Name:      name,
		Status:    status,
		CreatedBy: user,
		CreatedAt: now,
		UpdatedBy: &user,
		UpdatedAt: &now,
	}
	if err := s.repo.Create(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *fuelTypeService) UpdateFuelType(id int64, name string, status bool, user string) (*FuelType, error) {
	f, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, nil
	}

	now := time.Now().UTC()
	f.Name = name
	f.Status = status
	f.UpdatedBy = &user
	f.UpdatedAt = &now

	if err := s.repo.Update(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *fuelTypeService) DeleteFuelType(id int64) error {
	return s.repo.Delete(id)
}
