package vehicle

import "time"

type FuelCardService interface {
	ListAll() ([]FuelCard, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]FuelCard, int, error)
	FindByID(id int64) (*FuelCard, error)
	CreateFuelCard(cardNo string, fuelCompanyID int64, pinNumber string, vehicleID *int64, status bool, user string) (*FuelCard, error)
	UpdateFuelCard(id int64, cardNo string, fuelCompanyID int64, pinNumber string, vehicleID *int64, status bool, user string) (*FuelCard, error)
	DeleteFuelCard(id int64) error
}

type fuelCardService struct {
	repo FuelCardRepository
}

func NewFuelCardService(repo FuelCardRepository) FuelCardService {
	return &fuelCardService{repo: repo}
}

func (s *fuelCardService) ListAll() ([]FuelCard, error) {
	return s.repo.GetAll()
}

func (s *fuelCardService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]FuelCard, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	cards, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}
	return cards, total, nil
}

func (s *fuelCardService) FindByID(id int64) (*FuelCard, error) {
	return s.repo.GetByID(id)
}

func (s *fuelCardService) CreateFuelCard(cardNo string, fuelCompanyID int64, pinNumber string, vehicleID *int64, status bool, user string) (*FuelCard, error) {
	now := time.Now().UTC()
	fc := &FuelCard{
		CardNo:        cardNo,
		FuelCompanyID: fuelCompanyID,
		PinNumber:     pinNumber,
		VehicleID:     vehicleID,
		Status:        status,
		CreatedBy:     user,
		CreatedAt:     now,
		UpdatedBy:     &user,
		UpdatedAt:     &now,
	}
	if err := s.repo.Create(fc); err != nil {
		return nil, err
	}
	return fc, nil
}

func (s *fuelCardService) UpdateFuelCard(id int64, cardNo string, fuelCompanyID int64, pinNumber string, vehicleID *int64, status bool, user string) (*FuelCard, error) {
	fc, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if fc == nil {
		return nil, nil
	}

	now := time.Now().UTC()
	fc.CardNo = cardNo
	fc.FuelCompanyID = fuelCompanyID
	fc.PinNumber = pinNumber
	fc.VehicleID = vehicleID
	fc.Status = status
	fc.UpdatedBy = &user
	fc.UpdatedAt = &now

	if err := s.repo.Update(fc); err != nil {
		return nil, err
	}
	return fc, nil
}

func (s *fuelCardService) DeleteFuelCard(id int64) error {
	return s.repo.Delete(id)
}
