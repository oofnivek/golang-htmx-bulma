package vehicle

import (
	"database/sql"
)

type FuelCompanyRepository interface {
	GetAll() ([]FuelCompany, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]FuelCompany, error)
	Count() (int, error)
	GetByID(id int64) (*FuelCompany, error)
	Create(c *FuelCompany) error
	Update(c *FuelCompany) error
	Delete(id int64) error
}

type mysqlFuelCompanyRepository struct {
	db *sql.DB
}

func NewFuelCompanyRepository(db *sql.DB) FuelCompanyRepository {
	return &mysqlFuelCompanyRepository{db: db}
}

func (r *mysqlFuelCompanyRepository) GetAll() ([]FuelCompany, error) {
	rows, err := r.db.Query("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM fuel_company ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []FuelCompany
	for rows.Next() {
		var c FuelCompany
		if err := rows.Scan(&c.ID, &c.Name, &c.Status, &c.CreatedBy, &c.CreatedAt, &c.UpdatedBy, &c.UpdatedAt); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, nil
}

func (r *mysqlFuelCompanyRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]FuelCompany, error) {
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

	query := "SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM fuel_company ORDER BY " + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []FuelCompany
	for rows.Next() {
		var c FuelCompany
		if err := rows.Scan(&c.ID, &c.Name, &c.Status, &c.CreatedBy, &c.CreatedAt, &c.UpdatedBy, &c.UpdatedAt); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, nil
}

func (r *mysqlFuelCompanyRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM fuel_company").Scan(&count)
	return count, err
}

func (r *mysqlFuelCompanyRepository) GetByID(id int64) (*FuelCompany, error) {
	row := r.db.QueryRow("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM fuel_company WHERE id = ?", id)
	var c FuelCompany
	err := row.Scan(&c.ID, &c.Name, &c.Status, &c.CreatedBy, &c.CreatedAt, &c.UpdatedBy, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *mysqlFuelCompanyRepository) Create(c *FuelCompany) error {
	res, err := r.db.Exec("INSERT INTO fuel_company (name, status, created_by, created_at, updated_by, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		c.Name, c.Status, c.CreatedBy, c.CreatedAt, c.UpdatedBy, c.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (r *mysqlFuelCompanyRepository) Update(c *FuelCompany) error {
	_, err := r.db.Exec("UPDATE fuel_company SET name = ?, status = ?, updated_by = ?, updated_at = ? WHERE id = ?",
		c.Name, c.Status, c.UpdatedBy, c.UpdatedAt, c.ID)
	return err
}

func (r *mysqlFuelCompanyRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM fuel_company WHERE id = ?", id)
	return err
}
