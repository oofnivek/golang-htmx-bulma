package vehicle

import (
	"database/sql"
	"fmt"
)

type RegionalInfoRepository interface {
	GetAll() ([]RegionalInfo, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]RegionalInfo, error)
	Count() (int, error)
	GetByID(postalCode string) (*RegionalInfo, error)
	Create(r *RegionalInfo) error
	Update(r *RegionalInfo) error
	Delete(postalCode string) error
}

type mysqlRegionalInfoRepository struct {
	db *sql.DB
}

func NewRegionalInfoRepository(db *sql.DB) RegionalInfoRepository {
	return &mysqlRegionalInfoRepository{db: db}
}

const regionalInfoColumns = `
	ri.postal_code, ri.region, ri.estate_id,
	COALESCE(e.estate_name, '') AS estate_name`

const regionalInfoFrom = `
	FROM regional_info ri
	LEFT JOIN estate e ON ri.estate_id = e.id`

func scanRegionalInfo(row interface {
	Scan(dest ...any) error
}) (*RegionalInfo, error) {
	var r RegionalInfo
	if err := row.Scan(&r.PostalCode, &r.Region, &r.EstateID, &r.EstateName); err != nil {
		return nil, err
	}
	return &r, nil
}

var regionalInfoSortableColumns = map[string]bool{
	"postal_code": true,
	"region":      true,
	"estate_id":   true,
}

func (r *mysqlRegionalInfoRepository) GetAll() ([]RegionalInfo, error) {
	rows, err := r.db.Query("SELECT" + regionalInfoColumns + regionalInfoFrom + " ORDER BY ri.postal_code ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []RegionalInfo
	for rows.Next() {
		item, err := scanRegionalInfo(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	return items, nil
}

func (r *mysqlRegionalInfoRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]RegionalInfo, error) {
	if !regionalInfoSortableColumns[sortBy] {
		sortBy = "postal_code"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	query := fmt.Sprintf("SELECT%s%s ORDER BY ri.%s %s LIMIT ? OFFSET ?", regionalInfoColumns, regionalInfoFrom, sortBy, sortOrder)
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []RegionalInfo
	for rows.Next() {
		item, err := scanRegionalInfo(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	return items, nil
}

func (r *mysqlRegionalInfoRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM regional_info").Scan(&count)
	return count, err
}

func (r *mysqlRegionalInfoRepository) GetByID(postalCode string) (*RegionalInfo, error) {
	query := "SELECT" + regionalInfoColumns + regionalInfoFrom + " WHERE ri.postal_code = ?"
	row := r.db.QueryRow(query, postalCode)
	item, err := scanRegionalInfo(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return item, err
}

func (r *mysqlRegionalInfoRepository) Create(ri *RegionalInfo) error {
	_, err := r.db.Exec(
		"INSERT INTO regional_info (postal_code, region, estate_id) VALUES (?, ?, ?)",
		ri.PostalCode, ri.Region, ri.EstateID,
	)
	return err
}

func (r *mysqlRegionalInfoRepository) Update(ri *RegionalInfo) error {
	_, err := r.db.Exec(
		"UPDATE regional_info SET region = ?, estate_id = ? WHERE postal_code = ?",
		ri.Region, ri.EstateID, ri.PostalCode,
	)
	return err
}

func (r *mysqlRegionalInfoRepository) Delete(postalCode string) error {
	_, err := r.db.Exec("DELETE FROM regional_info WHERE postal_code = ?", postalCode)
	return err
}
