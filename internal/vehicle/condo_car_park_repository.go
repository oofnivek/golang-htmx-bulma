package vehicle

import (
	"database/sql"
)

type CondoCarParkRepository interface {
	GetAll() ([]CondoCarPark, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]CondoCarPark, error)
	Count() (int, error)
	GetByID(id int64) (*CondoCarPark, error)
	Create(c *CondoCarPark) error
	Update(c *CondoCarPark) error
	Delete(id int64) error
}

type mysqlCondoCarParkRepository struct {
	db *sql.DB
}

func NewCondoCarParkRepository(db *sql.DB) CondoCarParkRepository {
	return &mysqlCondoCarParkRepository{db: db}
}

const condoCarParkColumns = `
	ccp.id, ccp.condo_id,
	COALESCE(c.name, '') AS condo_name,
	ccp.car_park_id,
	COALESCE(cp.name, '') AS car_park_name`

const condoCarParkFrom = `
	FROM condo_car_park ccp
	LEFT JOIN condo c ON ccp.condo_id = c.id
	LEFT JOIN car_park cp ON ccp.car_park_id = cp.id`

func scanCondoCarPark(row interface {
	Scan(dest ...any) error
}) (*CondoCarPark, error) {
	var c CondoCarPark
	err := row.Scan(&c.ID, &c.CondoID, &c.CondoName, &c.CarParkID, &c.CarParkName)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *mysqlCondoCarParkRepository) GetAll() ([]CondoCarPark, error) {
	query := "SELECT" + condoCarParkColumns + condoCarParkFrom + " ORDER BY ccp.id DESC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CondoCarPark
	for rows.Next() {
		c, err := scanCondoCarPark(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *c)
	}
	return items, nil
}

func (r *mysqlCondoCarParkRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]CondoCarPark, error) {
	sortableColumns := map[string]bool{
		"id":          true,
		"condo_id":    true,
		"car_park_id": true,
	}
	if !sortableColumns[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := "SELECT" + condoCarParkColumns + condoCarParkFrom +
		" ORDER BY ccp." + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CondoCarPark
	for rows.Next() {
		c, err := scanCondoCarPark(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *c)
	}
	return items, nil
}

func (r *mysqlCondoCarParkRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM condo_car_park").Scan(&count)
	return count, err
}

func (r *mysqlCondoCarParkRepository) GetByID(id int64) (*CondoCarPark, error) {
	query := "SELECT" + condoCarParkColumns + condoCarParkFrom + " WHERE ccp.id = ?"
	row := r.db.QueryRow(query, id)
	c, err := scanCondoCarPark(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func (r *mysqlCondoCarParkRepository) Create(c *CondoCarPark) error {
	res, err := r.db.Exec(
		`INSERT INTO condo_car_park (condo_id, car_park_id) VALUES (?, ?)`,
		c.CondoID, c.CarParkID,
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

func (r *mysqlCondoCarParkRepository) Update(c *CondoCarPark) error {
	_, err := r.db.Exec(
		`UPDATE condo_car_park SET condo_id = ?, car_park_id = ? WHERE id = ?`,
		c.CondoID, c.CarParkID, c.ID,
	)
	return err
}

func (r *mysqlCondoCarParkRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM condo_car_park WHERE id = ?", id)
	return err
}
