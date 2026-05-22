package web

import (
	"golang-htmx-bulma/internal/vehicle"
)

type MockVehicleTypeService struct {
	ListAllFn   func() ([]vehicle.VehicleType, error)
	ListPagedFn func(page, pageSize int, sortBy, sortOrder string) ([]vehicle.VehicleType, int, error)
	FindByIDFn  func(id int64) (*vehicle.VehicleType, error)
	CreateFn    func(name string, status bool, user string) (*vehicle.VehicleType, error)
	UpdateFn    func(id int64, name string, status bool, user string) (*vehicle.VehicleType, error)
	DeleteFn    func(id int64) error
}

func (m *MockVehicleTypeService) ListAll() ([]vehicle.VehicleType, error) {
	return m.ListAllFn()
}

func (m *MockVehicleTypeService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]vehicle.VehicleType, int, error) {
	return m.ListPagedFn(page, pageSize, sortBy, sortOrder)
}

func (m *MockVehicleTypeService) FindByID(id int64) (*vehicle.VehicleType, error) {
	return m.FindByIDFn(id)
}

func (m *MockVehicleTypeService) Create(name string, status bool, user string) (*vehicle.VehicleType, error) {
	return m.CreateFn(name, status, user)
}

func (m *MockVehicleTypeService) Update(id int64, name string, status bool, user string) (*vehicle.VehicleType, error) {
	return m.UpdateFn(id, name, status, user)
}

func (m *MockVehicleTypeService) Delete(id int64) error {
	return m.DeleteFn(id)
}
