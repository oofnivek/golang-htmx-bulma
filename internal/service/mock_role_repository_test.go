package service

import (
	"golang-htmx-bulma/internal/model"
)

type MockRoleRepository struct {
	GetAllFn   func() ([]model.Role, error)
	GetPagedFn func(limit, offset int, sortBy, sortOrder, search string) ([]model.Role, error)
	CountFn    func(search string) (int, error)
	GetByIDFn  func(id string) (*model.Role, error)
	CreateFn   func(role *model.Role) error
	UpdateFn   func(role *model.Role) error
	DeleteFn   func(id string) error
}

func (m *MockRoleRepository) GetAll() ([]model.Role, error) {
	return m.GetAllFn()
}

func (m *MockRoleRepository) GetPaged(limit, offset int, sortBy, sortOrder, search string) ([]model.Role, error) {
	return m.GetPagedFn(limit, offset, sortBy, sortOrder, search)
}

func (m *MockRoleRepository) Count(search string) (int, error) {
	return m.CountFn(search)
}

func (m *MockRoleRepository) GetByID(id string) (*model.Role, error) {
	return m.GetByIDFn(id)
}

func (m *MockRoleRepository) Create(role *model.Role) error {
	return m.CreateFn(role)
}

func (m *MockRoleRepository) Update(role *model.Role) error {
	return m.UpdateFn(role)
}

func (m *MockRoleRepository) Delete(id string) error {
	return m.DeleteFn(id)
}
