package vehicle

import "time"

type VehicleTypeService interface {
    ListAll() ([]VehicleType, error)
    ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleType, int, error)
    FindByID(id int64) (*VehicleType, error)
    Create(name string, status bool, user string) (*VehicleType, error)
    Update(id int64, name string, status bool, user string) (*VehicleType, error)
    Delete(id int64) error
}

type vehicleTypeService struct {
    repo VehicleTypeRepository
}

func NewVehicleTypeService(r VehicleTypeRepository) VehicleTypeService {
    return &vehicleTypeService{repo: r}
}

func (s *vehicleTypeService) ListAll() ([]VehicleType, error) {
    return s.repo.GetAll()
}

func (s *vehicleTypeService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleType, int, error) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 {
        pageSize = 10
    }
    offset := (page - 1) * pageSize
    list, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
    if err != nil {
        return nil, 0, err
    }
    total, err := s.repo.Count()
    if err != nil {
        return nil, 0, err
    }
    return list, total, nil
}

func (s *vehicleTypeService) FindByID(id int64) (*VehicleType, error) {
    return s.repo.GetByID(id)
}

func (s *vehicleTypeService) Create(name string, status bool, user string) (*VehicleType, error) {
    now := time.Now()
    vt := &VehicleType{
        Name:      name,
        Status:    status,
        CreatedBy: user,
        CreatedAt: now,
        UpdatedBy: &user,
        UpdatedAt: &now,
    }
    if err := s.repo.Create(vt); err != nil {
        return nil, err
    }
    return vt, nil
}

func (s *vehicleTypeService) Update(id int64, name string, status bool, user string) (*VehicleType, error) {
    vt, err := s.repo.GetByID(id)
    if err != nil {
        return nil, err
    }
    if vt == nil {
        return nil, nil
    }
    vt.Name = name
    vt.Status = status
    now := time.Now()
    vt.UpdatedBy = &user
    vt.UpdatedAt = &now
    if err := s.repo.Update(vt); err != nil {
        return nil, err
    }
    return vt, nil
}

func (s *vehicleTypeService) Delete(id int64) error {
    return s.repo.Delete(id)
}
