package vehicle

import (
	"time"
)

type CarParkService interface {
	ListAll() ([]CarPark, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CarPark, int, error)
	FindByID(id int64) (*CarPark, error)
	CreateCarPark(name string, description *string, postalCode, address string, latitude, longitude float64, carParkOwnerID int64, activeFrom, activeTo *time.Time, status bool, user string) (*CarPark, error)
	UpdateCarPark(id int64, name string, description *string, postalCode, address string, latitude, longitude float64, carParkOwnerID int64, activeFrom, activeTo *time.Time, status bool, user string) (*CarPark, error)
	DeleteCarPark(id int64) error
}

type carParkService struct {
	repo CarParkRepository
}

func NewCarParkService(repo CarParkRepository) CarParkService {
	return &carParkService{repo: repo}
}

func (s *carParkService) ListAll() ([]CarPark, error) {
	return s.repo.GetAll()
}

func (s *carParkService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CarPark, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	parks, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return parks, total, nil
}

func (s *carParkService) FindByID(id int64) (*CarPark, error) {
	return s.repo.GetByID(id)
}

func (s *carParkService) CreateCarPark(name string, description *string, postalCode, address string, latitude, longitude float64, carParkOwnerID int64, activeFrom, activeTo *time.Time, status bool, user string) (*CarPark, error) {
	now := time.Now().UTC()
	cp := &CarPark{
		Name:           name,
		Description:    description,
		PostalCode:     postalCode,
		Address:        address,
		Latitude:       latitude,
		Longitude:      longitude,
		CarParkOwnerID: carParkOwnerID,
		ActiveFrom:     activeFrom,
		ActiveTo:       activeTo,
		Status:         status,
		CreatedBy:      user,
		CreatedAt:      now,
		UpdatedBy:      &user,
		UpdatedAt:      &now,
	}
	if err := s.repo.Create(cp); err != nil {
		return nil, err
	}
	return cp, nil
}

func (s *carParkService) UpdateCarPark(id int64, name string, description *string, postalCode, address string, latitude, longitude float64, carParkOwnerID int64, activeFrom, activeTo *time.Time, status bool, user string) (*CarPark, error) {
	cp, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if cp == nil {
		return nil, nil
	}

	now := time.Now().UTC()
	cp.Name = name
	cp.Description = description
	cp.PostalCode = postalCode
	cp.Address = address
	cp.Latitude = latitude
	cp.Longitude = longitude
	cp.CarParkOwnerID = carParkOwnerID
	cp.ActiveFrom = activeFrom
	cp.ActiveTo = activeTo
	cp.Status = status
	cp.UpdatedBy = &user
	cp.UpdatedAt = &now

	if err := s.repo.Update(cp); err != nil {
		return nil, err
	}
	return cp, nil
}

func (s *carParkService) DeleteCarPark(id int64) error {
	return s.repo.Delete(id)
}
