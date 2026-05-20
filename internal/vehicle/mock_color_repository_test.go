package vehicle

type MockVehicleColorRepository struct {
	GetAllFn   func() ([]VehicleColor, error)
	GetPagedFn func(limit, offset int, sortBy, sortOrder string) ([]VehicleColor, error)
	CountFn    func() (int, error)
	GetByIDFn  func(id int64) (*VehicleColor, error)
	CreateFn   func(color *VehicleColor) error
	UpdateFn   func(color *VehicleColor) error
	DeleteFn   func(id int64) error
}

func (m *MockVehicleColorRepository) GetAll() ([]VehicleColor, error) {
	return m.GetAllFn()
}

func (m *MockVehicleColorRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]VehicleColor, error) {
	return m.GetPagedFn(limit, offset, sortBy, sortOrder)
}

func (m *MockVehicleColorRepository) Count() (int, error) {
	return m.CountFn()
}

func (m *MockVehicleColorRepository) GetByID(id int64) (*VehicleColor, error) {
	return m.GetByIDFn(id)
}

func (m *MockVehicleColorRepository) Create(color *VehicleColor) error {
	return m.CreateFn(color)
}

func (m *MockVehicleColorRepository) Update(color *VehicleColor) error {
	return m.UpdateFn(color)
}

func (m *MockVehicleColorRepository) Delete(id int64) error {
	return m.DeleteFn(id)
}
