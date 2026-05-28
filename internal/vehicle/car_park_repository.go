package vehicle

import (
	"database/sql"
)

type CarParkRepository interface {
	GetAll() ([]CarPark, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]CarPark, error)
	Count() (int, error)
	GetByID(id int64) (*CarPark, error)
	Create(cp *CarPark) error
	Update(cp *CarPark) error
	Delete(id int64) error
}

type mysqlCarParkRepository struct {
	db *sql.DB
}

func NewCarParkRepository(db *sql.DB) CarParkRepository {
	return &mysqlCarParkRepository{db: db}
}

const carParkColumns = `
	cp.id, cp.name, cp.description, cp.postal_code, cp.address,
	cp.latitude, cp.longitude, cp.car_park_owner_id,
	COALESCE(cpo.name, '') AS car_park_owner_name,
	cp.active_from, cp.active_to, cp.status,
	cp.created_by, cp.created_at, cp.updated_by, cp.updated_at`

const carParkFrom = `
	FROM car_park cp
	LEFT JOIN car_park_owner cpo ON cp.car_park_owner_id = cpo.id`

func scanCarPark(row interface {
	Scan(dest ...any) error
}) (*CarPark, error) {
	var cp CarPark
	err := row.Scan(
		&cp.ID, &cp.Name, &cp.Description, &cp.PostalCode, &cp.Address,
		&cp.Latitude, &cp.Longitude, &cp.CarParkOwnerID, &cp.CarParkOwnerName,
		&cp.ActiveFrom, &cp.ActiveTo, &cp.Status,
		&cp.CreatedBy, &cp.CreatedAt, &cp.UpdatedBy, &cp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &cp, nil
}

func (r *mysqlCarParkRepository) GetAll() ([]CarPark, error) {
	query := "SELECT" + carParkColumns + carParkFrom + " ORDER BY cp.id DESC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parks []CarPark
	for rows.Next() {
		cp, err := scanCarPark(rows)
		if err != nil {
			return nil, err
		}
		parks = append(parks, *cp)
	}
	return parks, nil
}

func (r *mysqlCarParkRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]CarPark, error) {
	sortableColumns := map[string]bool{
		"id":          true,
		"name":        true,
		"postal_code": true,
		"status":      true,
		"updated_by":  true,
		"updated_at":  true,
	}
	if !sortableColumns[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := "SELECT" + carParkColumns + carParkFrom +
		" ORDER BY cp." + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parks []CarPark
	for rows.Next() {
		cp, err := scanCarPark(rows)
		if err != nil {
			return nil, err
		}
		parks = append(parks, *cp)
	}
	return parks, nil
}

func (r *mysqlCarParkRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM car_park").Scan(&count)
	return count, err
}

func (r *mysqlCarParkRepository) GetByID(id int64) (*CarPark, error) {
	query := "SELECT" + carParkColumns + carParkFrom + " WHERE cp.id = ?"
	row := r.db.QueryRow(query, id)
	cp, err := scanCarPark(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return cp, err
}

func (r *mysqlCarParkRepository) Create(cp *CarPark) error {
	res, err := r.db.Exec(
		`INSERT INTO car_park
			(name, description, postal_code, address, latitude, longitude,
			 car_park_owner_id, active_from, active_to, status,
			 created_by, created_at, updated_by, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		cp.Name, cp.Description, cp.PostalCode, cp.Address, cp.Latitude, cp.Longitude,
		cp.CarParkOwnerID, cp.ActiveFrom, cp.ActiveTo, cp.Status,
		cp.CreatedBy, cp.CreatedAt, cp.UpdatedBy, cp.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	cp.ID = id
	return nil
}

func (r *mysqlCarParkRepository) Update(cp *CarPark) error {
	_, err := r.db.Exec(
		`UPDATE car_park SET
			name = ?, description = ?, postal_code = ?, address = ?,
			latitude = ?, longitude = ?, car_park_owner_id = ?,
			active_from = ?, active_to = ?, status = ?,
			updated_by = ?, updated_at = ?
		 WHERE id = ?`,
		cp.Name, cp.Description, cp.PostalCode, cp.Address,
		cp.Latitude, cp.Longitude, cp.CarParkOwnerID,
		cp.ActiveFrom, cp.ActiveTo, cp.Status,
		cp.UpdatedBy, cp.UpdatedAt, cp.ID,
	)
	return err
}

func (r *mysqlCarParkRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM car_park WHERE id = ?", id)
	return err
}
