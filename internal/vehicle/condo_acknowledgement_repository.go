package vehicle

import (
	"database/sql"
	"fmt"
)

type CondoAcknowledgementRepository interface {
	GetAll() ([]CondoAcknowledgement, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]CondoAcknowledgement, error)
	Count() (int, error)
	GetByID(id int64) (*CondoAcknowledgement, error)
	Create(a *CondoAcknowledgement) error
	Delete(id int64) error
}

type mysqlCondoAcknowledgementRepository struct {
	db *sql.DB
}

func NewCondoAcknowledgementRepository(db *sql.DB) CondoAcknowledgementRepository {
	return &mysqlCondoAcknowledgementRepository{db: db}
}

const condoAcknowledgementColumns = "user_id, created_at, deleted_at"
const condoAcknowledgementFrom = "FROM condo_acknowledgement WHERE deleted_at IS NULL"

func scanCondoAcknowledgement(s interface{ Scan(...any) error }) (*CondoAcknowledgement, error) {
	var a CondoAcknowledgement
	if err := s.Scan(&a.UserID, &a.CreatedAt, &a.DeletedAt); err != nil {
		return nil, err
	}
	return &a, nil
}

var condoAcknowledgementSortableColumns = map[string]bool{
	"user_id":    true,
	"created_at": true,
}

func (r *mysqlCondoAcknowledgementRepository) GetAll() ([]CondoAcknowledgement, error) {
	rows, err := r.db.Query("SELECT " + condoAcknowledgementColumns + " " + condoAcknowledgementFrom + " ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CondoAcknowledgement
	for rows.Next() {
		a, err := scanCondoAcknowledgement(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *a)
	}
	return items, nil
}

func (r *mysqlCondoAcknowledgementRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]CondoAcknowledgement, error) {
	if !condoAcknowledgementSortableColumns[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := fmt.Sprintf("SELECT %s %s ORDER BY %s %s LIMIT ? OFFSET ?", condoAcknowledgementColumns, condoAcknowledgementFrom, sortBy, sortOrder)
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CondoAcknowledgement
	for rows.Next() {
		a, err := scanCondoAcknowledgement(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *a)
	}
	return items, nil
}

func (r *mysqlCondoAcknowledgementRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) " + condoAcknowledgementFrom).Scan(&count)
	return count, err
}

func (r *mysqlCondoAcknowledgementRepository) GetByID(id int64) (*CondoAcknowledgement, error) {
	row := r.db.QueryRow("SELECT "+condoAcknowledgementColumns+" "+condoAcknowledgementFrom+" AND user_id = ?", id)
	a, err := scanCondoAcknowledgement(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return a, err
}

func (r *mysqlCondoAcknowledgementRepository) Create(a *CondoAcknowledgement) error {
	_, err := r.db.Exec("INSERT INTO condo_acknowledgement (user_id, created_at) VALUES (?, ?)", a.UserID, a.CreatedAt)
	return err
}

func (r *mysqlCondoAcknowledgementRepository) Delete(id int64) error {
	_, err := r.db.Exec("UPDATE condo_acknowledgement SET deleted_at = NOW(3) WHERE user_id = ?", id)
	return err
}
