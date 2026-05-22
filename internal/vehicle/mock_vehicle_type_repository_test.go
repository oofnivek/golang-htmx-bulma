package vehicle

import (
    "fmt"
    "sort"
    "time"
)

type mockVehicleTypeRepository struct {
    data map[int64]*VehicleType
    nextID int64
    // error injection flags
    createErr  bool
    updateErr  bool
    getAllErr   bool
    getPagedErr bool
    countErr   bool
    getByIDErr bool
}

func newMockRepo() *mockVehicleTypeRepository {
    return &mockVehicleTypeRepository{data: make(map[int64]*VehicleType), nextID: 1}
}

func (m *mockVehicleTypeRepository) GetAll() ([]VehicleType, error) {
    if m.getAllErr {
        return nil, fmt.Errorf("GetAll error")
    }
    var list []VehicleType
    for _, vt := range m.data {
        list = append(list, *vt)
    }
    return list, nil
}

func (m *mockVehicleTypeRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]VehicleType, error) {
    if m.getPagedErr {
        return nil, fmt.Errorf("GetPaged error")
    }
    // ignore sorting for mock, just paginate based on insertion order
    var all []VehicleType
    for _, vt := range m.data {
        all = append(all, *vt)
    }
    // sort by ID to ensure deterministic ordering
    sort.Slice(all, func(i, j int) bool { return all[i].ID < all[j].ID })
    if offset > len(all) {
        return []VehicleType{}, nil
    }
    end := offset + limit
    if end > len(all) {
        end = len(all)
    }
    return all[offset:end], nil
}

func (m *mockVehicleTypeRepository) Count() (int, error) {
    if m.countErr {
        return 0, fmt.Errorf("Count error")
    }
    return len(m.data), nil
}

func (m *mockVehicleTypeRepository) GetByID(id int64) (*VehicleType, error) {
    if m.getByIDErr {
        return nil, fmt.Errorf("GetByID error")
    }
    if vt, ok := m.data[id]; ok {
        return vt, nil
    }
    return nil, nil
}

func (m *mockVehicleTypeRepository) Create(vt *VehicleType) error {
    if m.createErr {
        return fmt.Errorf("Create error")
    }
    vt.ID = m.nextID
    m.nextID++
    now := time.Now()
    vt.CreatedAt = now
    vt.UpdatedAt = &now
    m.data[vt.ID] = vt
    return nil
}

func (m *mockVehicleTypeRepository) Update(vt *VehicleType) error {
    if m.updateErr {
        return fmt.Errorf("Update error")
    }
    if _, ok := m.data[vt.ID]; !ok {
        return nil
    }
    now := time.Now()
    vt.UpdatedAt = &now
    m.data[vt.ID] = vt
    return nil
}

func (m *mockVehicleTypeRepository) Delete(id int64) error {
    delete(m.data, id)
    return nil
}
