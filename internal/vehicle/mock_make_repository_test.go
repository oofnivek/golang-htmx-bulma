package vehicle

type MockVehicleMakeRepository struct {
	GetAllFn   func() ([]VehicleMake, error)
	GetPagedFn func(limit, offset int, sortBy, sortOrder string) ([]VehicleMake, error)
	CountFn    func() (int, error)
	GetByIDFn  func(id int64) (*VehicleMake, error)
	CreateFn   func(make *VehicleMake) error
	UpdateFn   func(make *VehicleMake) error
	DeleteFn   func(id int64) error
}

func (m *MockVehicleMakeRepository) GetAll() ([]VehicleMake, error) {
	return m.GetAllFn()
}

func (m *MockVehicleMakeRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]VehicleMake, error) {
	return m.GetPagedFn(limit, offset, sortBy, sortOrder)
}

func (m *MockVehicleMakeRepository) Count() (int, error) {
	return m.CountFn()
}

func (m *MockVehicleMakeRepository) GetByID(id int64) (*VehicleMake, error) {
	return m.GetByIDFn(id)
}

func (m *MockVehicleMakeRepository) Create(make *VehicleMake) error {
	return m.CreateFn(make)
}

func (m *MockVehicleMakeRepository) Update(make *VehicleMake) error {
	return m.UpdateFn(make)
}

func (m *MockVehicleMakeRepository) Delete(id int64) error {
	return m.DeleteFn(id)
}
