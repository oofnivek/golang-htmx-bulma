package vehicle

import (
	"time"

	"golang-htmx-bulma/internal/pkg/status"
)

type CondoService interface {
	ListAll() ([]Condo, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]Condo, int, error)
	FindByID(id int64) (*Condo, error)
	CreateCondo(name string, st status.Status, mcstNumber string, mcstEmail string, address string, user string) (*Condo, error)
	UpdateCondo(id int64, name string, st status.Status, mcstNumber string, mcstEmail string, address string, user string) (*Condo, error)
	DeleteCondo(id int64) error
}

type condoService struct {
	repo CondoRepository
}

func NewCondoService(repo CondoRepository) CondoService {
	return &condoService{repo: repo}
}

func (s *condoService) ListAll() ([]Condo, error) {
	return s.repo.GetAll()
}

func (s *condoService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]Condo, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	condos, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}
	return condos, total, nil
}

func (s *condoService) FindByID(id int64) (*Condo, error) {
	return s.repo.GetByID(id)
}

func (s *condoService) CreateCondo(name string, st status.Status, mcstNumber string, mcstEmail string, address string, user string) (*Condo, error) {
	now := time.Now().UTC()
	c := &Condo{
		Name:       name,
		Status:     st,
		McstNumber: mcstNumber,
		McstEmail:  mcstEmail,
		Address:    address,
		CreatedBy:  user,
		CreatedAt:  now,
		UpdatedBy:  &user,
		UpdatedAt:  &now,
	}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *condoService) UpdateCondo(id int64, name string, st status.Status, mcstNumber string, mcstEmail string, address string, user string) (*Condo, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}

	now := time.Now().UTC()
	c.Name = name
	c.Status = st
	c.McstNumber = mcstNumber
	c.McstEmail = mcstEmail
	c.Address = address
	c.UpdatedBy = &user
	c.UpdatedAt = &now

	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *condoService) DeleteCondo(id int64) error {
	return s.repo.Delete(id)
}
