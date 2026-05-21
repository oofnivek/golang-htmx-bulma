package user

type MockRoleRepository struct {
	GetAllFn   func() ([]Role, error)
	GetPagedFn func(limit, offset int, sortBy, sortOrder, search string) ([]Role, error)
	CountFn    func(search string) (int, error)
	GetByIDFn  func(id string) (*Role, error)
	CreateFn   func(role *Role) error
	UpdateFn   func(role *Role) error
	DeleteFn   func(id string) error
}

func (m *MockRoleRepository) GetAll() ([]Role, error) {
	return m.GetAllFn()
}

func (m *MockRoleRepository) GetPaged(limit, offset int, sortBy, sortOrder, search string) ([]Role, error) {
	return m.GetPagedFn(limit, offset, sortBy, sortOrder, search)
}

func (m *MockRoleRepository) Count(search string) (int, error) {
	return m.CountFn(search)
}

func (m *MockRoleRepository) GetByID(id string) (*Role, error) {
	return m.GetByIDFn(id)
}

func (m *MockRoleRepository) Create(role *Role) error {
	return m.CreateFn(role)
}

func (m *MockRoleRepository) Update(role *Role) error {
	return m.UpdateFn(role)
}

func (m *MockRoleRepository) Delete(id string) error {
	return m.DeleteFn(id)
}
