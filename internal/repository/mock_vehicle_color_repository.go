package repository

import (
	"golang-htmx-bulma/internal/model"
)

type MockVehicleColorRepository struct {
	GetAllFn   func() ([]model.VehicleColor, error)
	GetByIDFn  func(id int64) (*model.VehicleColor, error)
	CreateFn   func(color *model.VehicleColor) error
	UpdateFn   func(color *model.VehicleColor) error
	DeleteFn   func(id int64) error
}

func (m *MockVehicleColorRepository) GetAll() ([]model.VehicleColor, error) {
	return m.GetAllFn()
}

func (m *MockVehicleColorRepository) GetByID(id int64) (*model.VehicleColor, error) {
	return m.GetByIDFn(id)
}

func (m *MockVehicleColorRepository) Create(color *model.VehicleColor) error {
	return m.CreateFn(color)
}

func (m *MockVehicleColorRepository) Update(color *model.VehicleColor) error {
	return m.UpdateFn(color)
}

func (m *MockVehicleColorRepository) Delete(id int64) error {
	return m.DeleteFn(id)
}
