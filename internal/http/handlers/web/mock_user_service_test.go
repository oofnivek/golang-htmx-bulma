package web

import (
	"golang-htmx-bulma/internal/user"
)

type MockUserService struct {
	ListAllFn    func() ([]user.User, error)
	ListPagedFn  func(page, pageSize int, sortBy, sortOrder, search string) ([]user.User, int, error)
	GetByEmailFn func(email string) (*user.User, error)
	CreateUserFn func(u *user.User) error
	UpdateUserFn func(u *user.User) error
	DeleteUserFn func(email string) error
}

func (m *MockUserService) ListAll() ([]user.User, error) {
	return m.ListAllFn()
}

func (m *MockUserService) ListPaged(page, pageSize int, sortBy, sortOrder, search string) ([]user.User, int, error) {
	return m.ListPagedFn(page, pageSize, sortBy, sortOrder, search)
}

func (m *MockUserService) GetByEmail(email string) (*user.User, error) {
	return m.GetByEmailFn(email)
}

func (m *MockUserService) CreateUser(u *user.User) error {
	return m.CreateUserFn(u)
}

func (m *MockUserService) UpdateUser(u *user.User) error {
	return m.UpdateUserFn(u)
}

func (m *MockUserService) DeleteUser(email string) error {
	return m.DeleteUserFn(email)
}
