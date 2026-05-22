package vehicle

import "database/sql"

type VehicleTypeRepository interface {
    GetAll() ([]VehicleType, error)
    GetPaged(limit, offset int, sortBy, sortOrder string) ([]VehicleType, error)
    Count() (int, error)
    GetByID(id int64) (*VehicleType, error)
    Create(vt *VehicleType) error
    Update(vt *VehicleType) error
    Delete(id int64) error
}

type mysqlVehicleTypeRepository struct { db *sql.DB }

func NewVehicleTypeRepository(db *sql.DB) VehicleTypeRepository { return &mysqlVehicleTypeRepository{db: db} }

func (r *mysqlVehicleTypeRepository) GetAll() ([]VehicleType, error) {
    rows, err := r.db.Query("SELECT id, name, status, created_by, created_at, updated_by, updated_at, old_id FROM vehicle_type ORDER BY id DESC")
    if err != nil { return nil, err }
    defer rows.Close()
    var types []VehicleType
    for rows.Next() {
        var vt VehicleType
        if err := rows.Scan(&vt.ID, &vt.Name, &vt.Status, &vt.CreatedBy, &vt.CreatedAt, &vt.UpdatedBy, &vt.UpdatedAt, &vt.OldID); err != nil {
            return nil, err
        }
        types = append(types, vt)
    }
    return types, nil
}

func (r *mysqlVehicleTypeRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]VehicleType, error) {
    // whitelist columns
    valid := map[string]bool{"id": true, "name": true, "status": true, "updated_by": true, "updated_at": true}
    if !valid[sortBy] { sortBy = "id" }
    if sortOrder != "asc" && sortOrder != "desc" { sortOrder = "desc" }
    query := "SELECT id, name, status, created_by, created_at, updated_by, updated_at, old_id FROM vehicle_type ORDER BY " + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
    rows, err := r.db.Query(query, limit, offset)
    if err != nil { return nil, err }
    defer rows.Close()
    var types []VehicleType
    for rows.Next() {
        var vt VehicleType
        if err := rows.Scan(&vt.ID, &vt.Name, &vt.Status, &vt.CreatedBy, &vt.CreatedAt, &vt.UpdatedBy, &vt.UpdatedAt, &vt.OldID); err != nil {
            return nil, err
        }
        types = append(types, vt)
    }
    return types, nil
}

func (r *mysqlVehicleTypeRepository) Count() (int, error) {
    var cnt int
    err := r.db.QueryRow("SELECT COUNT(*) FROM vehicle_type").Scan(&cnt)
    return cnt, err
}

func (r *mysqlVehicleTypeRepository) GetByID(id int64) (*VehicleType, error) {
    row := r.db.QueryRow("SELECT id, name, status, created_by, created_at, updated_by, updated_at, old_id FROM vehicle_type WHERE id = ?", id)
    var vt VehicleType
    err := row.Scan(&vt.ID, &vt.Name, &vt.Status, &vt.CreatedBy, &vt.CreatedAt, &vt.UpdatedBy, &vt.UpdatedAt, &vt.OldID)
    if err == sql.ErrNoRows { return nil, nil }
    if err != nil { return nil, err }
    return &vt, nil
}

func (r *mysqlVehicleTypeRepository) Create(vt *VehicleType) error {
    res, err := r.db.Exec("INSERT INTO vehicle_type (name, status, created_by, created_at, updated_by, updated_at, old_id) VALUES (?,?,?,?,?,?,?)",
        vt.Name, vt.Status, vt.CreatedBy, vt.CreatedAt, vt.UpdatedBy, vt.UpdatedAt, vt.OldID)
    if err != nil { return err }
    id, err := res.LastInsertId()
    if err != nil { return err }
    vt.ID = id
    return nil
}

func (r *mysqlVehicleTypeRepository) Update(vt *VehicleType) error {
    _, err := r.db.Exec("UPDATE vehicle_type SET name = ?, status = ?, updated_by = ?, updated_at = ?, old_id = ? WHERE id = ?",
        vt.Name, vt.Status, vt.UpdatedBy, vt.UpdatedAt, vt.OldID, vt.ID)
    return err
}

func (r *mysqlVehicleTypeRepository) Delete(id int64) error {
    _, err := r.db.Exec("DELETE FROM vehicle_type WHERE id = ?", id)
    return err
}
