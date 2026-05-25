package vehicle

import "database/sql"

type FuelTypeRepository interface {
	GetAll() ([]FuelType, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]FuelType, error)
	Count() (int, error)
	GetByID(id int64) (*FuelType, error)
	Create(f *FuelType) error
	Update(f *FuelType) error
	Delete(id int64) error
}

type mysqlFuelTypeRepository struct {
	db *sql.DB
}

func NewFuelTypeRepository(db *sql.DB) FuelTypeRepository {
	return &mysqlFuelTypeRepository{db: db}
}

func (r *mysqlFuelTypeRepository) GetAll() ([]FuelType, error) {
	rows, err := r.db.Query("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM fuel_type ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fuels []FuelType
	for rows.Next() {
		var f FuelType
		if err := rows.Scan(&f.ID, &f.Name, &f.Status, &f.CreatedBy, &f.CreatedAt, &f.UpdatedBy, &f.UpdatedAt); err != nil {
			return nil, err
		}
		fuels = append(fuels, f)
	}
	return fuels, nil
}

func (r *mysqlFuelTypeRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]FuelType, error) {
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

	query := "SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM fuel_type ORDER BY " + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fuels []FuelType
	for rows.Next() {
		var f FuelType
		if err := rows.Scan(&f.ID, &f.Name, &f.Status, &f.CreatedBy, &f.CreatedAt, &f.UpdatedBy, &f.UpdatedAt); err != nil {
			return nil, err
		}
		fuels = append(fuels, f)
	}
	return fuels, nil
}

func (r *mysqlFuelTypeRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM fuel_type").Scan(&count)
	return count, err
}

func (r *mysqlFuelTypeRepository) GetByID(id int64) (*FuelType, error) {
	row := r.db.QueryRow("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM fuel_type WHERE id = ?", id)
	var f FuelType
	err := row.Scan(&f.ID, &f.Name, &f.Status, &f.CreatedBy, &f.CreatedAt, &f.UpdatedBy, &f.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *mysqlFuelTypeRepository) Create(f *FuelType) error {
	res, err := r.db.Exec(
		"INSERT INTO fuel_type (name, status, created_by, created_at, updated_by, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		f.Name, f.Status, f.CreatedBy, f.CreatedAt, f.UpdatedBy, f.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	f.ID = id
	return nil
}

func (r *mysqlFuelTypeRepository) Update(f *FuelType) error {
	_, err := r.db.Exec(
		"UPDATE fuel_type SET name = ?, status = ?, updated_by = ?, updated_at = ? WHERE id = ?",
		f.Name, f.Status, f.UpdatedBy, f.UpdatedAt, f.ID,
	)
	return err
}

func (r *mysqlFuelTypeRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM fuel_type WHERE id = ?", id)
	return err
}
