package vehicle

import (
	"time"
)

type CarParkOwnerService interface {
	ListAll() ([]CarParkOwner, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CarParkOwner, int, error)
	FindByID(id int64) (*CarParkOwner, error)
	CreateOwner(name string, status bool, user string) (*CarParkOwner, error)
	UpdateOwner(id int64, name string, status bool, user string) (*CarParkOwner, error)
	DeleteOwner(id int64) error
}

type carParkOwnerService struct {
	repo CarParkOwnerRepository
}

func NewCarParkOwnerService(repo CarParkOwnerRepository) CarParkOwnerService {
	return &carParkOwnerService{repo: repo}
}

func (s *carParkOwnerService) ListAll() ([]CarParkOwner, error) {
	return s.repo.GetAll()
}

func (s *carParkOwnerService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CarParkOwner, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	owners, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return owners, total, nil
}

func (s *carParkOwnerService) FindByID(id int64) (*CarParkOwner, error) {
	return s.repo.GetByID(id)
}

func (s *carParkOwnerService) CreateOwner(name string, status bool, user string) (*CarParkOwner, error) {
	now := time.Now().UTC()
	owner := &CarParkOwner{
		Name:      name,
		Status:    status,
		CreatedBy: user,
		CreatedAt: now,
		UpdatedBy: &user,
		UpdatedAt: &now,
	}
	if err := s.repo.Create(owner); err != nil {
		return nil, err
	}
	return owner, nil
}

func (s *carParkOwnerService) UpdateOwner(id int64, name string, status bool, user string) (*CarParkOwner, error) {
	owner, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if owner == nil {
		return nil, nil
	}

	owner.Name = name
	owner.Status = status
	now := time.Now().UTC()
	owner.UpdatedBy = &user
	owner.UpdatedAt = &now

	if err := s.repo.Update(owner); err != nil {
		return nil, err
	}
	return owner, nil
}

func (s *carParkOwnerService) DeleteOwner(id int64) error {
	return s.repo.Delete(id)
}
