package web

import (
	"golang-htmx-bulma/internal/model"
)

type MockVehicleMakeService struct {
	ListAllFn       func() ([]model.VehicleMake, error)
	FindByIDFn      func(id int64) (*model.VehicleMake, error)
	CreateMakeFn    func(name string, status bool, user string) (*model.VehicleMake, error)
	UpdateMakeFn    func(id int64, name string, status bool, user string) (*model.VehicleMake, error)
	DeleteMakeFn    func(id int64) error
}

func (m *MockVehicleMakeService) ListAll() ([]model.VehicleMake, error) {
	return m.ListAllFn()
}

func (m *MockVehicleMakeService) FindByID(id int64) (*model.VehicleMake, error) {
	return m.FindByIDFn(id)
}

func (m *MockVehicleMakeService) CreateMake(name string, status bool, user string) (*model.VehicleMake, error) {
	return m.CreateMakeFn(name, status, user)
}

func (m *MockVehicleMakeService) UpdateMake(id int64, name string, status bool, user string) (*model.VehicleMake, error) {
	return m.UpdateMakeFn(id, name, status, user)
}

func (m *MockVehicleMakeService) DeleteMake(id int64) error {
	return m.DeleteMakeFn(id)
}
