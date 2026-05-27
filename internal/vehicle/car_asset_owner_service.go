package vehicle

import (
	"time"
)

type CarAssetOwnerService interface {
	ListAll() ([]CarAssetOwner, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CarAssetOwner, int, error)
	FindByID(id int64) (*CarAssetOwner, error)
	CreateOwner(name string, status bool, user string) (*CarAssetOwner, error)
	UpdateOwner(id int64, name string, status bool, user string) (*CarAssetOwner, error)
	DeleteOwner(id int64) error
}

type carAssetOwnerService struct {
	repo CarAssetOwnerRepository
}

func NewCarAssetOwnerService(repo CarAssetOwnerRepository) CarAssetOwnerService {
	return &carAssetOwnerService{repo: repo}
}

func (s *carAssetOwnerService) ListAll() ([]CarAssetOwner, error) {
	return s.repo.GetAll()
}

func (s *carAssetOwnerService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CarAssetOwner, int, error) {
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

func (s *carAssetOwnerService) FindByID(id int64) (*CarAssetOwner, error) {
	return s.repo.GetByID(id)
}

func (s *carAssetOwnerService) CreateOwner(name string, status bool, user string) (*CarAssetOwner, error) {
	now := time.Now().UTC()
	owner := &CarAssetOwner{
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

func (s *carAssetOwnerService) UpdateOwner(id int64, name string, status bool, user string) (*CarAssetOwner, error) {
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

func (s *carAssetOwnerService) DeleteOwner(id int64) error {
	return s.repo.Delete(id)
}
