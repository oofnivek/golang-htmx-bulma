package vehicle

import (
	"database/sql"
)

type VehicleStatusRepository interface {
	GetAll() ([]VehicleStatus, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]VehicleStatus, error)
	Count() (int, error)
	GetByID(id int64) (*VehicleStatus, error)
	Create(s *VehicleStatus) error
	Update(s *VehicleStatus) error
	Delete(id int64) error
}

type mysqlVehicleStatusRepository struct {
	db *sql.DB
}

func NewVehicleStatusRepository(db *sql.DB) VehicleStatusRepository {
	return &mysqlVehicleStatusRepository{db: db}
}

func (r *mysqlVehicleStatusRepository) GetAll() ([]VehicleStatus, error) {
	rows, err := r.db.Query("SELECT id, is_active, substatus FROM vehicle_status ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []VehicleStatus
	for rows.Next() {
		var s VehicleStatus
		if err := rows.Scan(&s.ID, &s.IsActive, &s.Substatus); err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}
	return statuses, nil
}

func (r *mysqlVehicleStatusRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]VehicleStatus, error) {
	validColumns := map[string]bool{
		"id":        true,
		"is_active": true,
		"substatus": true,
	}
	if !validColumns[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := "SELECT id, is_active, substatus FROM vehicle_status ORDER BY " + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []VehicleStatus
	for rows.Next() {
		var s VehicleStatus
		if err := rows.Scan(&s.ID, &s.IsActive, &s.Substatus); err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}
	return statuses, nil
}

func (r *mysqlVehicleStatusRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM vehicle_status").Scan(&count)
	return count, err
}

func (r *mysqlVehicleStatusRepository) GetByID(id int64) (*VehicleStatus, error) {
	row := r.db.QueryRow("SELECT id, is_active, substatus FROM vehicle_status WHERE id = ?", id)
	var s VehicleStatus
	err := row.Scan(&s.ID, &s.IsActive, &s.Substatus)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *mysqlVehicleStatusRepository) Create(s *VehicleStatus) error {
	res, err := r.db.Exec("INSERT INTO vehicle_status (is_active, substatus) VALUES (?, ?)", s.IsActive, s.Substatus)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	s.ID = id
	return nil
}

func (r *mysqlVehicleStatusRepository) Update(s *VehicleStatus) error {
	_, err := r.db.Exec("UPDATE vehicle_status SET is_active = ?, substatus = ? WHERE id = ?", s.IsActive, s.Substatus, s.ID)
	return err
}

func (r *mysqlVehicleStatusRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM vehicle_status WHERE id = ?", id)
	return err
}
