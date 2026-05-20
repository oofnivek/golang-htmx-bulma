package vehicle

import (
	"database/sql"
)

type VehicleMakeRepository interface {
	GetAll() ([]VehicleMake, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]VehicleMake, error)
	Count() (int, error)
	GetByID(id int64) (*VehicleMake, error)
	Create(make *VehicleMake) error
	Update(make *VehicleMake) error
	Delete(id int64) error
}

type mysqlVehicleMakeRepository struct {
	db *sql.DB
}

func NewVehicleMakeRepository(db *sql.DB) VehicleMakeRepository {
	return &mysqlVehicleMakeRepository{db: db}
}

func (r *mysqlVehicleMakeRepository) GetAll() ([]VehicleMake, error) {
	rows, err := r.db.Query("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_make ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var makes []VehicleMake
	for rows.Next() {
		var m VehicleMake
		err := rows.Scan(&m.ID, &m.Name, &m.Status, &m.CreatedBy, &m.CreatedAt, &m.UpdatedBy, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		makes = append(makes, m)
	}
	return makes, nil
}

func (r *mysqlVehicleMakeRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]VehicleMake, error) {
	// Whitelist sorting columns to prevent SQL injection
	validColumns := map[string]bool{
		"id":         true,
		"name":       true,
		"status":     true,
		"updated_by": true,
		"updated_at": true,
	}
	if !validColumns[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := "SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_make ORDER BY " + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var makes []VehicleMake
	for rows.Next() {
		var m VehicleMake
		err := rows.Scan(&m.ID, &m.Name, &m.Status, &m.CreatedBy, &m.CreatedAt, &m.UpdatedBy, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		makes = append(makes, m)
	}
	return makes, nil
}

func (r *mysqlVehicleMakeRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM vehicle_make").Scan(&count)
	return count, err
}

func (r *mysqlVehicleMakeRepository) GetByID(id int64) (*VehicleMake, error) {
	row := r.db.QueryRow("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_make WHERE id = ?", id)
	var m VehicleMake
	err := row.Scan(&m.ID, &m.Name, &m.Status, &m.CreatedBy, &m.CreatedAt, &m.UpdatedBy, &m.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *mysqlVehicleMakeRepository) Create(m *VehicleMake) error {
	res, err := r.db.Exec("INSERT INTO vehicle_make (name, status, created_by, created_at, updated_by, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		m.Name, m.Status, m.CreatedBy, m.CreatedAt, m.UpdatedBy, m.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	m.ID = id
	return nil
}

func (r *mysqlVehicleMakeRepository) Update(m *VehicleMake) error {
	_, err := r.db.Exec("UPDATE vehicle_make SET name = ?, status = ?, updated_by = ?, updated_at = ? WHERE id = ?",
		m.Name, m.Status, m.UpdatedBy, m.UpdatedAt, m.ID)
	return err
}

func (r *mysqlVehicleMakeRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM vehicle_make WHERE id = ?", id)
	return err
}
