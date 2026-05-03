package service

import (
	"golang-htmx-bulma/internal/model"
	"golang-htmx-bulma/internal/repository"
)

type RoleService interface {
	ListAll() ([]model.Role, error)
	ListPaged(page, pageSize int, sortBy, sortOrder, search string) ([]model.Role, int, error)
	FindByID(id string) (*model.Role, error)
	CreateRole(id, name string) (*model.Role, error)
	UpdateRole(id, name string) (*model.Role, error)
	DeleteRole(id string) error
}

type roleService struct {
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) ListAll() ([]model.Role, error) {
	return s.repo.GetAll()
}

func (s *roleService) ListPaged(page, pageSize int, sortBy, sortOrder, search string) ([]model.Role, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	roles, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder, search)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(search)
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (s *roleService) FindByID(id string) (*model.Role, error) {
	return s.repo.GetByID(id)
}

func (s *roleService) CreateRole(id, name string) (*model.Role, error) {
	role := &model.Role{
		ID:   id,
		Name: name,
	}
	err := s.repo.Create(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *roleService) UpdateRole(id, name string) (*model.Role, error) {
	role, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	role.Name = name

	err = s.repo.Update(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *roleService) DeleteRole(id string) error {
	return s.repo.Delete(id)
}
