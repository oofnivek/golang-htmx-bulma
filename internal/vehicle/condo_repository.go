package vehicle

import (
	"database/sql"
	"fmt"
)

type CondoRepository interface {
	GetAll() ([]Condo, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]Condo, error)
	Count() (int, error)
	GetByID(id int64) (*Condo, error)
	Create(c *Condo) error
	Update(c *Condo) error
	Delete(id int64) error
}

type mysqlCondoRepository struct {
	db *sql.DB
}

func NewCondoRepository(db *sql.DB) CondoRepository {
	return &mysqlCondoRepository{db: db}
}

const condoColumns = "id, name, status, mcst_number, mcst_email, address, created_by, created_at, updated_by, updated_at"
const condoFrom = "FROM condo"

func scanCondo(row interface{ Scan(dest ...any) error }) (*Condo, error) {
	var c Condo
	err := row.Scan(
		&c.ID, &c.Name, &c.Status, &c.McstNumber, &c.McstEmail, &c.Address,
		&c.CreatedBy, &c.CreatedAt, &c.UpdatedBy, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

var condoSortableColumns = map[string]bool{
	"id":          true,
	"name":        true,
	"status":      true,
	"mcst_number": true,
	"address":     true,
	"created_at":  true,
	"updated_at":  true,
}

func (r *mysqlCondoRepository) GetAll() ([]Condo, error) {
	rows, err := r.db.Query("SELECT " + condoColumns + " " + condoFrom + " ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var condos []Condo
	for rows.Next() {
		c, err := scanCondo(rows)
		if err != nil {
			return nil, err
		}
		condos = append(condos, *c)
	}
	return condos, nil
}

func (r *mysqlCondoRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]Condo, error) {
	if !condoSortableColumns[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := fmt.Sprintf("SELECT %s %s ORDER BY %s %s LIMIT ? OFFSET ?", condoColumns, condoFrom, sortBy, sortOrder)
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var condos []Condo
	for rows.Next() {
		c, err := scanCondo(rows)
		if err != nil {
			return nil, err
		}
		condos = append(condos, *c)
	}
	return condos, nil
}

func (r *mysqlCondoRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) " + condoFrom).Scan(&count)
	return count, err
}

func (r *mysqlCondoRepository) GetByID(id int64) (*Condo, error) {
	row := r.db.QueryRow("SELECT "+condoColumns+" "+condoFrom+" WHERE id = ?", id)
	c, err := scanCondo(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func (r *mysqlCondoRepository) Create(c *Condo) error {
	res, err := r.db.Exec(
		`INSERT INTO condo (name, status, mcst_number, mcst_email, address, created_by, created_at, updated_by, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.Name, c.Status, c.McstNumber, c.McstEmail, c.Address,
		c.CreatedBy, c.CreatedAt, c.UpdatedBy, c.UpdatedAt,
	)
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

func (r *mysqlCondoRepository) Update(c *Condo) error {
	_, err := r.db.Exec(
		`UPDATE condo SET name = ?, status = ?, mcst_number = ?, mcst_email = ?, address = ?, updated_by = ?, updated_at = ?
		 WHERE id = ?`,
		c.Name, c.Status, c.McstNumber, c.McstEmail, c.Address, c.UpdatedBy, c.UpdatedAt, c.ID,
	)
	return err
}

func (r *mysqlCondoRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM condo WHERE id = ?", id)
	return err
}
