package vehicle

import "time"

type CondoAcknowledgementService interface {
	ListAll() ([]CondoAcknowledgement, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CondoAcknowledgement, int, error)
	FindByID(id int64) (*CondoAcknowledgement, error)
	CreateCondoAcknowledgement(userID int64) (*CondoAcknowledgement, error)
	DeleteCondoAcknowledgement(id int64) error
}

type condoAcknowledgementService struct {
	repo CondoAcknowledgementRepository
}

func NewCondoAcknowledgementService(repo CondoAcknowledgementRepository) CondoAcknowledgementService {
	return &condoAcknowledgementService{repo: repo}
}

func (s *condoAcknowledgementService) ListAll() ([]CondoAcknowledgement, error) {
	return s.repo.GetAll()
}

func (s *condoAcknowledgementService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CondoAcknowledgement, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	items, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *condoAcknowledgementService) FindByID(id int64) (*CondoAcknowledgement, error) {
	return s.repo.GetByID(id)
}

func (s *condoAcknowledgementService) CreateCondoAcknowledgement(userID int64) (*CondoAcknowledgement, error) {
	a := &CondoAcknowledgement{
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *condoAcknowledgementService) DeleteCondoAcknowledgement(id int64) error {
	return s.repo.Delete(id)
}
