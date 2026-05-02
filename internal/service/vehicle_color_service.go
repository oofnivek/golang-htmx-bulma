package service

import (
	"golang-htmx-bulma/internal/model"
	"golang-htmx-bulma/internal/repository"
	"time"
)

type VehicleColorService interface {
	ListAll() ([]model.VehicleColor, error)
	FindByID(id int64) (*model.VehicleColor, error)
	CreateColor(name string, status bool, user string) (*model.VehicleColor, error)
	UpdateColor(id int64, name string, status bool, user string) (*model.VehicleColor, error)
	DeleteColor(id int64) error
}

type vehicleColorService struct {
	repo repository.VehicleColorRepository
}

func NewVehicleColorService(repo repository.VehicleColorRepository) VehicleColorService {
	return &vehicleColorService{repo: repo}
}

func (s *vehicleColorService) ListAll() ([]model.VehicleColor, error) {
	return s.repo.GetAll()
}

func (s *vehicleColorService) FindByID(id int64) (*model.VehicleColor, error) {
	return s.repo.GetByID(id)
}

func (s *vehicleColorService) CreateColor(name string, status bool, user string) (*model.VehicleColor, error) {
	now := time.Now()
	color := &model.VehicleColor{
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

func (s *vehicleColorService) UpdateColor(id int64, name string, status bool, user string) (*model.VehicleColor, error) {
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
