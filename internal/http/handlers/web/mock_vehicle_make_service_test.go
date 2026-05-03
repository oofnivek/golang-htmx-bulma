package web

import (
	"golang-htmx-bulma/internal/model"
)

type MockVehicleMakeService struct {
	ListAllFn       func() ([]model.VehicleMake, error)
	ListPagedFn     func(page, pageSize int, sortBy, sortOrder, search string) ([]model.VehicleMake, int, error)
	FindByIDFn      func(id int64) (*model.VehicleMake, error)
	CreateMakeFn    func(name string, status bool, user string) (*model.VehicleMake, error)
	UpdateMakeFn    func(id int64, name string, status bool, user string) (*model.VehicleMake, error)
	DeleteMakeFn    func(id int64) error
}

func (m *MockVehicleMakeService) ListAll() ([]model.VehicleMake, error) {
	return m.ListAllFn()
}

func (m *MockVehicleMakeService) ListPaged(page, pageSize int, sortBy, sortOrder, search string) ([]model.VehicleMake, int, error) {
	return m.ListPagedFn(page, pageSize, sortBy, sortOrder, search)
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
