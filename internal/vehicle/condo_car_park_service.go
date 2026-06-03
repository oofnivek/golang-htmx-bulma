package vehicle

type CondoCarParkService interface {
	ListAll() ([]CondoCarPark, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CondoCarPark, int, error)
	FindByID(id int64) (*CondoCarPark, error)
	CreateCondoCarPark(condoID, carParkID int64) (*CondoCarPark, error)
	UpdateCondoCarPark(id, condoID, carParkID int64) (*CondoCarPark, error)
	DeleteCondoCarPark(id int64) error
}

type condoCarParkService struct {
	repo CondoCarParkRepository
}

func NewCondoCarParkService(repo CondoCarParkRepository) CondoCarParkService {
	return &condoCarParkService{repo: repo}
}

func (s *condoCarParkService) ListAll() ([]CondoCarPark, error) {
	return s.repo.GetAll()
}

func (s *condoCarParkService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CondoCarPark, int, error) {
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

func (s *condoCarParkService) FindByID(id int64) (*CondoCarPark, error) {
	return s.repo.GetByID(id)
}

func (s *condoCarParkService) CreateCondoCarPark(condoID, carParkID int64) (*CondoCarPark, error) {
	c := &CondoCarPark{
		CondoID:   condoID,
		CarParkID: carParkID,
	}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *condoCarParkService) UpdateCondoCarPark(id, condoID, carParkID int64) (*CondoCarPark, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}

	c.CondoID = condoID
	c.CarParkID = carParkID

	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *condoCarParkService) DeleteCondoCarPark(id int64) error {
	return s.repo.Delete(id)
}
