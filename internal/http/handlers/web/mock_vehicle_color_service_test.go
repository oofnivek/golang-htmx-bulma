package web

import (
	"golang-htmx-bulma/internal/vehicle"
)

type MockVehicleColorService struct {
	ListAllFn      func() ([]vehicle.VehicleColor, error)
	ListPagedFn    func(page, pageSize int, sortBy, sortOrder string) ([]vehicle.VehicleColor, int, error)
	FindByIDFn     func(id int64) (*vehicle.VehicleColor, error)
	CreateColorFn  func(name string, status bool, user string) (*vehicle.VehicleColor, error)
	UpdateColorFn  func(id int64, name string, status bool, user string) (*vehicle.VehicleColor, error)
	DeleteColorFn  func(id int64) error
}

func (m *MockVehicleColorService) ListAll() ([]vehicle.VehicleColor, error) {
	return m.ListAllFn()
}

func (m *MockVehicleColorService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]vehicle.VehicleColor, int, error) {
	return m.ListPagedFn(page, pageSize, sortBy, sortOrder)
}

func (m *MockVehicleColorService) FindByID(id int64) (*vehicle.VehicleColor, error) {
	return m.FindByIDFn(id)
}

func (m *MockVehicleColorService) CreateColor(name string, status bool, user string) (*vehicle.VehicleColor, error) {
	return m.CreateColorFn(name, status, user)
}

func (m *MockVehicleColorService) UpdateColor(id int64, name string, status bool, user string) (*vehicle.VehicleColor, error) {
	return m.UpdateColorFn(id, name, status, user)
}

func (m *MockVehicleColorService) DeleteColor(id int64) error {
	return m.DeleteColorFn(id)
}
