package service

import (
	"golang-htmx-bulma/internal/model"
	"golang-htmx-bulma/internal/repository"
	"time"
)

type UserService interface {
	ListAll() ([]model.User, error)
	ListPaged(page, pageSize int, sortBy, sortOrder, search string) ([]model.User, int, error)
	GetByEmail(email string) (*model.User, error)
	CreateUser(u *model.User) error
	UpdateUser(u *model.User) error
	DeleteUser(email string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) ListAll() ([]model.User, error) {
	return s.repo.GetAll()
}

func (s *userService) ListPaged(page, pageSize int, sortBy, sortOrder, search string) ([]model.User, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	users, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder, search)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(search)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *userService) GetByEmail(email string) (*model.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *userService) CreateUser(u *model.User) error {
	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now
	return s.repo.Create(u)
}

func (s *userService) UpdateUser(u *model.User) error {
	u.UpdatedAt = time.Now().UTC()
	return s.repo.Update(u)
}

func (s *userService) DeleteUser(email string) error {
	return s.repo.Delete(email)
}
