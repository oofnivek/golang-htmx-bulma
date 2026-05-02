package repository

import (
	"database/sql"
	"golang-htmx-bulma/internal/model"
)

type VehicleColorRepository interface {
	GetAll() ([]model.VehicleColor, error)
	GetPaged(limit, offset int, sortBy, sortOrder, search string) ([]model.VehicleColor, error)
	Count(search string) (int, error)
	GetByID(id int64) (*model.VehicleColor, error)
	Create(color *model.VehicleColor) error
	Update(color *model.VehicleColor) error
	Delete(id int64) error
}

type mysqlVehicleColorRepository struct {
	db *sql.DB
}

func NewVehicleColorRepository(db *sql.DB) VehicleColorRepository {
	return &mysqlVehicleColorRepository{db: db}
}

func (r *mysqlVehicleColorRepository) GetAll() ([]model.VehicleColor, error) {
	rows, err := r.db.Query("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_color ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var colors []model.VehicleColor
	for rows.Next() {
		var c model.VehicleColor
		err := rows.Scan(&c.ID, &c.Name, &c.Status, &c.CreatedBy, &c.CreatedAt, &c.UpdatedBy, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		colors = append(colors, c)
	}
	return colors, nil
}

func (r *mysqlVehicleColorRepository) GetPaged(limit, offset int, sortBy, sortOrder, search string) ([]model.VehicleColor, error) {
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

	query := "SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_color"
	var args []interface{}

	if len(search) >= 2 {
		// Use FULLTEXT MATCH...AGAINST with n-gram index for fast partial matching.
		// Requires: ALTER TABLE vehicle_color ADD FULLTEXT INDEX idx_color_name_ngram (name) WITH PARSER ngram;
		query += " WHERE MATCH(name) AGAINST(? IN BOOLEAN MODE)"
		args = append(args, search)
	}

	query += " ORDER BY " + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var colors []model.VehicleColor
	for rows.Next() {
		var c model.VehicleColor
		err := rows.Scan(&c.ID, &c.Name, &c.Status, &c.CreatedBy, &c.CreatedAt, &c.UpdatedBy, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		colors = append(colors, c)
	}
	return colors, nil
}

func (r *mysqlVehicleColorRepository) Count(search string) (int, error) {
	query := "SELECT COUNT(*) FROM vehicle_color"
	var args []interface{}
	if search != "" {
		// Use FULLTEXT MATCH...AGAINST with n-gram index for fast partial matching.
		// Requires: ALTER TABLE vehicle_color ADD FULLTEXT INDEX idx_color_name_ngram (name) WITH PARSER ngram;
		query += " WHERE MATCH(name) AGAINST(? IN BOOLEAN MODE)"
		args = append(args, search)
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

func (r *mysqlVehicleColorRepository) GetByID(id int64) (*model.VehicleColor, error) {
	row := r.db.QueryRow("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM vehicle_color WHERE id = ?", id)
	var c model.VehicleColor
	err := row.Scan(&c.ID, &c.Name, &c.Status, &c.CreatedBy, &c.CreatedAt, &c.UpdatedBy, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *mysqlVehicleColorRepository) Create(c *model.VehicleColor) error {
	res, err := r.db.Exec("INSERT INTO vehicle_color (name, status, created_by, created_at, updated_by, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
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

func (r *mysqlVehicleColorRepository) Update(c *model.VehicleColor) error {
	_, err := r.db.Exec("UPDATE vehicle_color SET name = ?, status = ?, updated_by = ?, updated_at = ? WHERE id = ?",
		c.Name, c.Status, c.UpdatedBy, c.UpdatedAt, c.ID)
	return err
}

func (r *mysqlVehicleColorRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM vehicle_color WHERE id = ?", id)
	return err
}
