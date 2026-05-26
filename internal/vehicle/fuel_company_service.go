package vehicle

import "time"

type FuelCompanyService interface {
	ListAll() ([]FuelCompany, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]FuelCompany, int, error)
	FindByID(id int64) (*FuelCompany, error)
	CreateFuelCompany(name string, status bool, user string) (*FuelCompany, error)
	UpdateFuelCompany(id int64, name string, status bool, user string) (*FuelCompany, error)
	DeleteFuelCompany(id int64) error
}

type fuelCompanyService struct {
	repo FuelCompanyRepository
}

func NewFuelCompanyService(repo FuelCompanyRepository) FuelCompanyService {
	return &fuelCompanyService{repo: repo}
}

func (s *fuelCompanyService) ListAll() ([]FuelCompany, error) {
	return s.repo.GetAll()
}

func (s *fuelCompanyService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]FuelCompany, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	companies, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return companies, total, nil
}

func (s *fuelCompanyService) FindByID(id int64) (*FuelCompany, error) {
	return s.repo.GetByID(id)
}

func (s *fuelCompanyService) CreateFuelCompany(name string, status bool, user string) (*FuelCompany, error) {
	now := time.Now().UTC()
	c := &FuelCompany{
		Name:      name,
		Status:    status,
		CreatedBy: user,
		CreatedAt: now,
		UpdatedBy: &user,
		UpdatedAt: &now,
	}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *fuelCompanyService) UpdateFuelCompany(id int64, name string, status bool, user string) (*FuelCompany, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}

	c.Name = name
	c.Status = status
	now := time.Now().UTC()
	c.UpdatedBy = &user
	c.UpdatedAt = &now

	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *fuelCompanyService) DeleteFuelCompany(id int64) error {
	return s.repo.Delete(id)
}
