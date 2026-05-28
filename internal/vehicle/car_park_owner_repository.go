package vehicle

import (
	"database/sql"
)

type CarParkOwnerRepository interface {
	GetAll() ([]CarParkOwner, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]CarParkOwner, error)
	Count() (int, error)
	GetByID(id int64) (*CarParkOwner, error)
	Create(owner *CarParkOwner) error
	Update(owner *CarParkOwner) error
	Delete(id int64) error
}

type mysqlCarParkOwnerRepository struct {
	db *sql.DB
}

func NewCarParkOwnerRepository(db *sql.DB) CarParkOwnerRepository {
	return &mysqlCarParkOwnerRepository{db: db}
}

func (r *mysqlCarParkOwnerRepository) GetAll() ([]CarParkOwner, error) {
	rows, err := r.db.Query("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM car_park_owner ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var owners []CarParkOwner
	for rows.Next() {
		var o CarParkOwner
		if err := rows.Scan(&o.ID, &o.Name, &o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedBy, &o.UpdatedAt); err != nil {
			return nil, err
		}
		owners = append(owners, o)
	}
	return owners, nil
}

func (r *mysqlCarParkOwnerRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]CarParkOwner, error) {
	sortableColumns := map[string]bool{
		"id":         true,
		"name":       true,
		"status":     true,
		"updated_by": true,
		"updated_at": true,
	}
	if !sortableColumns[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := "SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM car_park_owner ORDER BY " + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var owners []CarParkOwner
	for rows.Next() {
		var o CarParkOwner
		if err := rows.Scan(&o.ID, &o.Name, &o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedBy, &o.UpdatedAt); err != nil {
			return nil, err
		}
		owners = append(owners, o)
	}
	return owners, nil
}

func (r *mysqlCarParkOwnerRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM car_park_owner").Scan(&count)
	return count, err
}

func (r *mysqlCarParkOwnerRepository) GetByID(id int64) (*CarParkOwner, error) {
	row := r.db.QueryRow("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM car_park_owner WHERE id = ?", id)
	var o CarParkOwner
	err := row.Scan(&o.ID, &o.Name, &o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedBy, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *mysqlCarParkOwnerRepository) Create(o *CarParkOwner) error {
	res, err := r.db.Exec("INSERT INTO car_park_owner (name, status, created_by, created_at, updated_by, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		o.Name, o.Status, o.CreatedBy, o.CreatedAt, o.UpdatedBy, o.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	o.ID = id
	return nil
}

func (r *mysqlCarParkOwnerRepository) Update(o *CarParkOwner) error {
	_, err := r.db.Exec("UPDATE car_park_owner SET name = ?, status = ?, updated_by = ?, updated_at = ? WHERE id = ?",
		o.Name, o.Status, o.UpdatedBy, o.UpdatedAt, o.ID)
	return err
}

func (r *mysqlCarParkOwnerRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM car_park_owner WHERE id = ?", id)
	return err
}
