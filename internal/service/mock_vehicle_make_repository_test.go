package service

import (
	"golang-htmx-bulma/internal/model"
)

type MockVehicleMakeRepository struct {
	GetAllFn   func() ([]model.VehicleMake, error)
	GetPagedFn func(limit, offset int) ([]model.VehicleMake, error)
	CountFn    func() (int, error)
	GetByIDFn  func(id int64) (*model.VehicleMake, error)
	CreateFn   func(make *model.VehicleMake) error
	UpdateFn   func(make *model.VehicleMake) error
	DeleteFn   func(id int64) error
}

func (m *MockVehicleMakeRepository) GetAll() ([]model.VehicleMake, error) {
	return m.GetAllFn()
}

func (m *MockVehicleMakeRepository) GetPaged(limit, offset int) ([]model.VehicleMake, error) {
	return m.GetPagedFn(limit, offset)
}

func (m *MockVehicleMakeRepository) Count() (int, error) {
	return m.CountFn()
}

func (m *MockVehicleMakeRepository) GetByID(id int64) (*model.VehicleMake, error) {
	return m.GetByIDFn(id)
}

func (m *MockVehicleMakeRepository) Create(make *model.VehicleMake) error {
	return m.CreateFn(make)
}

func (m *MockVehicleMakeRepository) Update(make *model.VehicleMake) error {
	return m.UpdateFn(make)
}

func (m *MockVehicleMakeRepository) Delete(id int64) error {
	return m.DeleteFn(id)
}
