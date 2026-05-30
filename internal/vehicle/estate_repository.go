package vehicle

import (
	"database/sql"
	"fmt"
)

type EstateRepository interface {
	GetAll() ([]Estate, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]Estate, error)
	Count() (int, error)
	GetByID(id int64) (*Estate, error)
	Create(e *Estate) error
	Update(e *Estate) error
	Delete(id int64) error
}

type mysqlEstateRepository struct {
	db *sql.DB
}

func NewEstateRepository(db *sql.DB) EstateRepository {
	return &mysqlEstateRepository{db: db}
}

const estateColumns = "id, estate_name"
const estateFrom = "FROM estate"

func scanEstate(s interface{ Scan(...any) error }) (*Estate, error) {
	var e Estate
	if err := s.Scan(&e.ID, &e.Name); err != nil {
		return nil, err
	}
	return &e, nil
}

var estateSortableColumns = map[string]bool{
	"id":          true,
	"estate_name": true,
}

func (r *mysqlEstateRepository) GetAll() ([]Estate, error) {
	rows, err := r.db.Query("SELECT " + estateColumns + " " + estateFrom + " ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var estates []Estate
	for rows.Next() {
		e, err := scanEstate(rows)
		if err != nil {
			return nil, err
		}
		estates = append(estates, *e)
	}
	return estates, nil
}

func (r *mysqlEstateRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]Estate, error) {
	if !estateSortableColumns[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := fmt.Sprintf("SELECT %s %s ORDER BY %s %s LIMIT ? OFFSET ?", estateColumns, estateFrom, sortBy, sortOrder)
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var estates []Estate
	for rows.Next() {
		e, err := scanEstate(rows)
		if err != nil {
			return nil, err
		}
		estates = append(estates, *e)
	}
	return estates, nil
}

func (r *mysqlEstateRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) " + estateFrom).Scan(&count)
	return count, err
}

func (r *mysqlEstateRepository) GetByID(id int64) (*Estate, error) {
	row := r.db.QueryRow("SELECT "+estateColumns+" "+estateFrom+" WHERE id = ?", id)
	e, err := scanEstate(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return e, err
}

func (r *mysqlEstateRepository) Create(e *Estate) error {
	res, err := r.db.Exec("INSERT INTO estate (estate_name) VALUES (?)", e.Name)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	e.ID = id
	return nil
}

func (r *mysqlEstateRepository) Update(e *Estate) error {
	_, err := r.db.Exec("UPDATE estate SET estate_name = ? WHERE id = ?", e.Name, e.ID)
	return err
}

func (r *mysqlEstateRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM estate WHERE id = ?", id)
	return err
}
