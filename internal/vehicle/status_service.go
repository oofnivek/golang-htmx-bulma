package vehicle

type VehicleStatusService interface {
	ListAll() ([]VehicleStatus, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleStatus, int, error)
	FindByID(id int64) (*VehicleStatus, error)
	CreateStatus(substatus string, isActive bool) (*VehicleStatus, error)
	UpdateStatus(id int64, substatus string, isActive bool) (*VehicleStatus, error)
	DeleteStatus(id int64) error
}

type vehicleStatusService struct {
	repo VehicleStatusRepository
}

func NewVehicleStatusService(repo VehicleStatusRepository) VehicleStatusService {
	return &vehicleStatusService{repo: repo}
}

func (s *vehicleStatusService) ListAll() ([]VehicleStatus, error) {
	return s.repo.GetAll()
}

func (s *vehicleStatusService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleStatus, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	statuses, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return statuses, total, nil
}

func (s *vehicleStatusService) FindByID(id int64) (*VehicleStatus, error) {
	return s.repo.GetByID(id)
}

func (s *vehicleStatusService) CreateStatus(substatus string, isActive bool) (*VehicleStatus, error) {
	vs := &VehicleStatus{
		Substatus: substatus,
		IsActive:  isActive,
	}
	if err := s.repo.Create(vs); err != nil {
		return nil, err
	}
	return vs, nil
}

func (s *vehicleStatusService) UpdateStatus(id int64, substatus string, isActive bool) (*VehicleStatus, error) {
	vs, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if vs == nil {
		return nil, nil
	}

	vs.Substatus = substatus
	vs.IsActive = isActive

	if err := s.repo.Update(vs); err != nil {
		return nil, err
	}
	return vs, nil
}

func (s *vehicleStatusService) DeleteStatus(id int64) error {
	return s.repo.Delete(id)
}
