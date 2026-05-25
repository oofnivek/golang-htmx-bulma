package vehicle

import (
	"time"
)

type VehicleMakeService interface {
	ListAll() ([]VehicleMake, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleMake, int, error)
	FindByID(id int64) (*VehicleMake, error)
	CreateMake(name string, status bool, user string) (*VehicleMake, error)
	UpdateMake(id int64, name string, status bool, user string) (*VehicleMake, error)
	DeleteMake(id int64) error
}

type vehicleMakeService struct {
	repo VehicleMakeRepository
}

func NewVehicleMakeService(repo VehicleMakeRepository) VehicleMakeService {
	return &vehicleMakeService{repo: repo}
}

func (s *vehicleMakeService) ListAll() ([]VehicleMake, error) {
	return s.repo.GetAll()
}

func (s *vehicleMakeService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleMake, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	makes, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return makes, total, nil
}

func (s *vehicleMakeService) FindByID(id int64) (*VehicleMake, error) {
	return s.repo.GetByID(id)
}

func (s *vehicleMakeService) CreateMake(name string, status bool, user string) (*VehicleMake, error) {
	now := time.Now().UTC()
	make := &VehicleMake{
		Name:      name,
		Status:    status,
		CreatedBy: user,
		CreatedAt: now,
		UpdatedBy: &user,
		UpdatedAt: &now,
	}
	err := s.repo.Create(make)
	if err != nil {
		return nil, err
	}
	return make, nil
}

func (s *vehicleMakeService) UpdateMake(id int64, name string, status bool, user string) (*VehicleMake, error) {
	make, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if make == nil {
		return nil, nil
	}

	make.Name = name
	make.Status = status
	now := time.Now().UTC()
	make.UpdatedBy = &user
	make.UpdatedAt = &now

	err = s.repo.Update(make)
	if err != nil {
		return nil, err
	}
	return make, nil
}

func (s *vehicleMakeService) DeleteMake(id int64) error {
	return s.repo.Delete(id)
}
