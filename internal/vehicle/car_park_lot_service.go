package vehicle

import "time"

type CarParkLotService interface {
	ListAll() ([]CarParkLot, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CarParkLot, int, error)
	FindByID(id int64) (*CarParkLot, error)
	CreateCarParkLot(carParkID int64, lotNumber, level string, status bool, user string) (*CarParkLot, error)
	UpdateCarParkLot(id, carParkID int64, lotNumber, level string, status bool, user string) (*CarParkLot, error)
	DeleteCarParkLot(id int64) error
}

type carParkLotService struct {
	repo CarParkLotRepository
}

func NewCarParkLotService(repo CarParkLotRepository) CarParkLotService {
	return &carParkLotService{repo: repo}
}

func (s *carParkLotService) ListAll() ([]CarParkLot, error) {
	return s.repo.GetAll()
}

func (s *carParkLotService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CarParkLot, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	lots, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return lots, total, nil
}

func (s *carParkLotService) FindByID(id int64) (*CarParkLot, error) {
	return s.repo.GetByID(id)
}

func (s *carParkLotService) CreateCarParkLot(carParkID int64, lotNumber, level string, status bool, user string) (*CarParkLot, error) {
	now := time.Now().UTC()
	l := &CarParkLot{
		CarParkID: carParkID,
		LotNumber: lotNumber,
		Level:     level,
		Status:    status,
		CreatedBy: user,
		CreatedAt: now,
		UpdatedBy: &user,
		UpdatedAt: &now,
	}
	if err := s.repo.Create(l); err != nil {
		return nil, err
	}
	return l, nil
}

func (s *carParkLotService) UpdateCarParkLot(id, carParkID int64, lotNumber, level string, status bool, user string) (*CarParkLot, error) {
	l, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if l == nil {
		return nil, nil
	}

	now := time.Now().UTC()
	l.CarParkID = carParkID
	l.LotNumber = lotNumber
	l.Level = level
	l.Status = status
	l.UpdatedBy = &user
	l.UpdatedAt = &now

	if err := s.repo.Update(l); err != nil {
		return nil, err
	}
	return l, nil
}

func (s *carParkLotService) DeleteCarParkLot(id int64) error {
	return s.repo.Delete(id)
}
