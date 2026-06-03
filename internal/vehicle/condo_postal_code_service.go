package vehicle

type CondoPostalCodeService interface {
	ListAll() ([]CondoPostalCode, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CondoPostalCode, int, error)
	FindByID(id int64) (*CondoPostalCode, error)
	CreateCondoPostalCode(condoID int64, postalCode string) (*CondoPostalCode, error)
	UpdateCondoPostalCode(id, condoID int64, postalCode string) (*CondoPostalCode, error)
	DeleteCondoPostalCode(id int64) error
}

type condoPostalCodeService struct {
	repo CondoPostalCodeRepository
}

func NewCondoPostalCodeService(repo CondoPostalCodeRepository) CondoPostalCodeService {
	return &condoPostalCodeService{repo: repo}
}

func (s *condoPostalCodeService) ListAll() ([]CondoPostalCode, error) {
	return s.repo.GetAll()
}

func (s *condoPostalCodeService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CondoPostalCode, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	codes, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}
	return codes, total, nil
}

func (s *condoPostalCodeService) FindByID(id int64) (*CondoPostalCode, error) {
	return s.repo.GetByID(id)
}

func (s *condoPostalCodeService) CreateCondoPostalCode(condoID int64, postalCode string) (*CondoPostalCode, error) {
	c := &CondoPostalCode{
		CondoID:    condoID,
		PostalCode: postalCode,
	}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *condoPostalCodeService) UpdateCondoPostalCode(id, condoID int64, postalCode string) (*CondoPostalCode, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}

	c.CondoID = condoID
	c.PostalCode = postalCode
	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *condoPostalCodeService) DeleteCondoPostalCode(id int64) error {
	return s.repo.Delete(id)
}
