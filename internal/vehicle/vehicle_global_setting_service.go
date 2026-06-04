package vehicle

import "time"

type VehicleGlobalSettingService interface {
	ListAll() ([]VehicleGlobalSetting, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleGlobalSetting, int, error)
	FindByID(id int64) (*VehicleGlobalSetting, error)
	CreateVehicleGlobalSetting(key, value string, remark, countryCode *string, user string) (*VehicleGlobalSetting, error)
	UpdateVehicleGlobalSetting(id int64, key, value string, remark, countryCode *string, user string) (*VehicleGlobalSetting, error)
	DeleteVehicleGlobalSetting(id int64) error
}

type vehicleGlobalSettingService struct {
	repo VehicleGlobalSettingRepository
}

func NewVehicleGlobalSettingService(repo VehicleGlobalSettingRepository) VehicleGlobalSettingService {
	return &vehicleGlobalSettingService{repo: repo}
}

func (s *vehicleGlobalSettingService) ListAll() ([]VehicleGlobalSetting, error) {
	return s.repo.GetAll()
}

func (s *vehicleGlobalSettingService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleGlobalSetting, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	settings, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}
	return settings, total, nil
}

func (s *vehicleGlobalSettingService) FindByID(id int64) (*VehicleGlobalSetting, error) {
	return s.repo.GetByID(id)
}

func (s *vehicleGlobalSettingService) CreateVehicleGlobalSetting(key, value string, remark, countryCode *string, user string) (*VehicleGlobalSetting, error) {
	now := time.Now().UTC()
	g := &VehicleGlobalSetting{
		Key:         key,
		Value:       value,
		Remark:      remark,
		CountryCode: countryCode,
		CreatedBy:   user,
		CreatedAt:   now,
		UpdatedBy:   &user,
		UpdatedAt:   &now,
	}
	if err := s.repo.Create(g); err != nil {
		return nil, err
	}
	return g, nil
}

func (s *vehicleGlobalSettingService) UpdateVehicleGlobalSetting(id int64, key, value string, remark, countryCode *string, user string) (*VehicleGlobalSetting, error) {
	g, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, nil
	}

	now := time.Now().UTC()
	g.Key = key
	g.Value = value
	g.Remark = remark
	g.CountryCode = countryCode
	g.UpdatedBy = &user
	g.UpdatedAt = &now

	if err := s.repo.Update(g); err != nil {
		return nil, err
	}
	return g, nil
}

func (s *vehicleGlobalSettingService) DeleteVehicleGlobalSetting(id int64) error {
	return s.repo.Delete(id)
}
