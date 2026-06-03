package vehicle

import (
	"database/sql"
	"fmt"
)

type CondoPostalCodeRepository interface {
	GetAll() ([]CondoPostalCode, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]CondoPostalCode, error)
	Count() (int, error)
	GetByID(id int64) (*CondoPostalCode, error)
	Create(c *CondoPostalCode) error
	Update(c *CondoPostalCode) error
	Delete(id int64) error
}

type mysqlCondoPostalCodeRepository struct {
	db *sql.DB
}

func NewCondoPostalCodeRepository(db *sql.DB) CondoPostalCodeRepository {
	return &mysqlCondoPostalCodeRepository{db: db}
}

const condoPostalCodeColumns = `
	cpc.id, cpc.condo_id,
	COALESCE(c.name, '') AS condo_name,
	cpc.postal_code`

const condoPostalCodeFrom = `
	FROM condo_postal_code cpc
	LEFT JOIN condo c ON cpc.condo_id = c.id`

func scanCondoPostalCode(row interface {
	Scan(dest ...any) error
}) (*CondoPostalCode, error) {
	var c CondoPostalCode
	err := row.Scan(&c.ID, &c.CondoID, &c.CondoName, &c.PostalCode)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

var condoPostalCodeSortableColumns = map[string]bool{
	"id":          true,
	"postal_code": true,
}

func (r *mysqlCondoPostalCodeRepository) GetAll() ([]CondoPostalCode, error) {
	query := "SELECT" + condoPostalCodeColumns + condoPostalCodeFrom + " ORDER BY cpc.id DESC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []CondoPostalCode
	for rows.Next() {
		c, err := scanCondoPostalCode(rows)
		if err != nil {
			return nil, err
		}
		codes = append(codes, *c)
	}
	return codes, nil
}

func (r *mysqlCondoPostalCodeRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]CondoPostalCode, error) {
	if !condoPostalCodeSortableColumns[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := fmt.Sprintf("SELECT%s%s ORDER BY cpc.%s %s LIMIT ? OFFSET ?",
		condoPostalCodeColumns, condoPostalCodeFrom, sortBy, sortOrder)
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []CondoPostalCode
	for rows.Next() {
		c, err := scanCondoPostalCode(rows)
		if err != nil {
			return nil, err
		}
		codes = append(codes, *c)
	}
	return codes, nil
}

func (r *mysqlCondoPostalCodeRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM condo_postal_code").Scan(&count)
	return count, err
}

func (r *mysqlCondoPostalCodeRepository) GetByID(id int64) (*CondoPostalCode, error) {
	query := "SELECT" + condoPostalCodeColumns + condoPostalCodeFrom + " WHERE cpc.id = ?"
	row := r.db.QueryRow(query, id)
	c, err := scanCondoPostalCode(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func (r *mysqlCondoPostalCodeRepository) Create(c *CondoPostalCode) error {
	res, err := r.db.Exec(
		"INSERT INTO condo_postal_code (condo_id, postal_code) VALUES (?, ?)",
		c.CondoID, c.PostalCode,
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

func (r *mysqlCondoPostalCodeRepository) Update(c *CondoPostalCode) error {
	_, err := r.db.Exec(
		"UPDATE condo_postal_code SET condo_id = ?, postal_code = ? WHERE id = ?",
		c.CondoID, c.PostalCode, c.ID,
	)
	return err
}

func (r *mysqlCondoPostalCodeRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM condo_postal_code WHERE id = ?", id)
	return err
}
