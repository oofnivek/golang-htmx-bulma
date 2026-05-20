package vehicle

import (
	"time"
)

type VehicleColorService interface {
	ListAll() ([]VehicleColor, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleColor, int, error)
	FindByID(id int64) (*VehicleColor, error)
	CreateColor(name string, status bool, user string) (*VehicleColor, error)
	UpdateColor(id int64, name string, status bool, user string) (*VehicleColor, error)
	DeleteColor(id int64) error
}

type vehicleColorService struct {
	repo VehicleColorRepository
}

func NewVehicleColorService(repo VehicleColorRepository) VehicleColorService {
	return &vehicleColorService{repo: repo}
}

func (s *vehicleColorService) ListAll() ([]VehicleColor, error) {
	return s.repo.GetAll()
}

func (s *vehicleColorService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleColor, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	colors, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return colors, total, nil
}

func (s *vehicleColorService) FindByID(id int64) (*VehicleColor, error) {
	return s.repo.GetByID(id)
}

func (s *vehicleColorService) CreateColor(name string, status bool, user string) (*VehicleColor, error) {
	now := time.Now()
	color := &VehicleColor{
		Name:      name,
		Status:    status,
		CreatedBy: user,
		CreatedAt: now,
		UpdatedBy: &user,
		UpdatedAt: &now,
	}
	err := s.repo.Create(color)
	if err != nil {
		return nil, err
	}
	return color, nil
}

func (s *vehicleColorService) UpdateColor(id int64, name string, status bool, user string) (*VehicleColor, error) {
	color, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if color == nil {
		return nil, nil
	}

	color.Name = name
	color.Status = status
	now := time.Now()
	color.UpdatedBy = &user
	color.UpdatedAt = &now

	err = s.repo.Update(color)
	if err != nil {
		return nil, err
	}
	return color, nil
}

func (s *vehicleColorService) DeleteColor(id int64) error {
	return s.repo.Delete(id)
}
