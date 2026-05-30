package vehicle

type EstateService interface {
	ListAll() ([]Estate, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]Estate, int, error)
	FindByID(id int64) (*Estate, error)
	CreateEstate(name string) (*Estate, error)
	UpdateEstate(id int64, name string) (*Estate, error)
	DeleteEstate(id int64) error
}

type estateService struct {
	repo EstateRepository
}

func NewEstateService(repo EstateRepository) EstateService {
	return &estateService{repo: repo}
}

func (s *estateService) ListAll() ([]Estate, error) {
	return s.repo.GetAll()
}

func (s *estateService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]Estate, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	estates, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}
	return estates, total, nil
}

func (s *estateService) FindByID(id int64) (*Estate, error) {
	return s.repo.GetByID(id)
}

func (s *estateService) CreateEstate(name string) (*Estate, error) {
	e := &Estate{Name: name}
	if err := s.repo.Create(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *estateService) UpdateEstate(id int64, name string) (*Estate, error) {
	e, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, nil
	}

	e.Name = name
	if err := s.repo.Update(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *estateService) DeleteEstate(id int64) error {
	return s.repo.Delete(id)
}
