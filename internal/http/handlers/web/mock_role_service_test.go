package web

import (
	"golang-htmx-bulma/internal/user"
)

type MockRoleService struct {
	ListAllFn    func() ([]user.Role, error)
	ListPagedFn  func(page, pageSize int, sortBy, sortOrder string) ([]user.Role, int, error)
	FindByIDFn   func(id string) (*user.Role, error)
	CreateRoleFn func(id, name string) (*user.Role, error)
	UpdateRoleFn func(id, name string) (*user.Role, error)
	DeleteRoleFn func(id string) error
}

func (m *MockRoleService) ListAll() ([]user.Role, error) {
	return m.ListAllFn()
}

func (m *MockRoleService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]user.Role, int, error) {
	return m.ListPagedFn(page, pageSize, sortBy, sortOrder)
}

func (m *MockRoleService) FindByID(id string) (*user.Role, error) {
	return m.FindByIDFn(id)
}

func (m *MockRoleService) CreateRole(id, name string) (*user.Role, error) {
	return m.CreateRoleFn(id, name)
}

func (m *MockRoleService) UpdateRole(id, name string) (*user.Role, error) {
	return m.UpdateRoleFn(id, name)
}

func (m *MockRoleService) DeleteRole(id string) error {
	return m.DeleteRoleFn(id)
}
