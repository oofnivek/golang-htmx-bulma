package web

import (
	"golang-htmx-bulma/internal/model"
)

type MockRoleService struct {
	ListAllFn       func() ([]model.Role, error)
	ListPagedFn     func(page, pageSize int, sortBy, sortOrder, search string) ([]model.Role, int, error)
	FindByIDFn      func(id string) (*model.Role, error)
	CreateRoleFn    func(id, name string) (*model.Role, error)
	UpdateRoleFn    func(id, name string) (*model.Role, error)
	DeleteRoleFn    func(id string) error
}

func (m *MockRoleService) ListAll() ([]model.Role, error) {
	return m.ListAllFn()
}

func (m *MockRoleService) ListPaged(page, pageSize int, sortBy, sortOrder, search string) ([]model.Role, int, error) {
	return m.ListPagedFn(page, pageSize, sortBy, sortOrder, search)
}

func (m *MockRoleService) FindByID(id string) (*model.Role, error) {
	return m.FindByIDFn(id)
}

func (m *MockRoleService) CreateRole(id, name string) (*model.Role, error) {
	return m.CreateRoleFn(id, name)
}

func (m *MockRoleService) UpdateRole(id, name string) (*model.Role, error) {
	return m.UpdateRoleFn(id, name)
}

func (m *MockRoleService) DeleteRole(id string) error {
	return m.DeleteRoleFn(id)
}
