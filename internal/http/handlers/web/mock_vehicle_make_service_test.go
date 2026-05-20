package web

import (
	"golang-htmx-bulma/internal/vehicle"
)

type MockVehicleMakeService struct {
	ListAllFn       func() ([]vehicle.VehicleMake, error)
	ListPagedFn     func(page, pageSize int, sortBy, sortOrder string) ([]vehicle.VehicleMake, int, error)
	FindByIDFn      func(id int64) (*vehicle.VehicleMake, error)
	CreateMakeFn    func(name string, status bool, user string) (*vehicle.VehicleMake, error)
	UpdateMakeFn    func(id int64, name string, status bool, user string) (*vehicle.VehicleMake, error)
	DeleteMakeFn    func(id int64) error
}

func (m *MockVehicleMakeService) ListAll() ([]vehicle.VehicleMake, error) {
	return m.ListAllFn()
}

func (m *MockVehicleMakeService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]vehicle.VehicleMake, int, error) {
	return m.ListPagedFn(page, pageSize, sortBy, sortOrder)
}

func (m *MockVehicleMakeService) FindByID(id int64) (*vehicle.VehicleMake, error) {
	return m.FindByIDFn(id)
}

func (m *MockVehicleMakeService) CreateMake(name string, status bool, user string) (*vehicle.VehicleMake, error) {
	return m.CreateMakeFn(name, status, user)
}

func (m *MockVehicleMakeService) UpdateMake(id int64, name string, status bool, user string) (*vehicle.VehicleMake, error) {
	return m.UpdateMakeFn(id, name, status, user)
}

func (m *MockVehicleMakeService) DeleteMake(id int64) error {
	return m.DeleteMakeFn(id)
}
