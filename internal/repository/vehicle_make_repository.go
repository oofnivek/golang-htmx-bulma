package repository

import (
	"database/sql"
	"golang-htmx-bulma/internal/model"
)

type VehicleMakeRepository interface {
	GetAll() ([]model.VehicleMake, error)
	GetByID(id int64) (*model.VehicleMake, error)
	Create(make *model.VehicleMake) error
	Update(make *model.VehicleMake) error
	Delete(id int64) error
}

type mysqlVehicleMakeRepository struct {
	db *sql.DB
}

func NewVehicleMakeRepository(db *sql.DB) VehicleMakeRepository {
	return &mysqlVehicleMakeRepository{db: db}
}

func (r *mysqlVehicleMakeRepository) GetAll() ([]model.VehicleMake, error) {
	rows, err := r.db.Query("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_make")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var makes []model.VehicleMake
	for rows.Next() {
		var m model.VehicleMake
		err := rows.Scan(&m.ID, &m.Name, &m.Status, &m.CreatedBy, &m.CreatedAt, &m.UpdatedBy, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		makes = append(makes, m)
	}
	return makes, nil
}

func (r *mysqlVehicleMakeRepository) GetByID(id int64) (*model.VehicleMake, error) {
	row := r.db.QueryRow("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_make WHERE id = ?", id)
	var m model.VehicleMake
	err := row.Scan(&m.ID, &m.Name, &m.Status, &m.CreatedBy, &m.CreatedAt, &m.UpdatedBy, &m.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *mysqlVehicleMakeRepository) Create(m *model.VehicleMake) error {
	res, err := r.db.Exec("INSERT INTO vehicle_make (name, status, created_by, created_at) VALUES (?, ?, ?, ?)",
		m.Name, m.Status, m.CreatedBy, m.CreatedAt)
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

func (r *mysqlVehicleMakeRepository) Update(m *model.VehicleMake) error {
	_, err := r.db.Exec("UPDATE vehicle_make SET name = ?, status = ?, updated_by = ?, updated_at = ? WHERE id = ?",
		m.Name, m.Status, m.UpdatedBy, m.UpdatedAt, m.ID)
	return err
}

func (r *mysqlVehicleMakeRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM vehicle_make WHERE id = ?", id)
	return err
}
