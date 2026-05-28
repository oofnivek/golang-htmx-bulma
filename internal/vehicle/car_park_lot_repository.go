package vehicle

import (
	"database/sql"
)

type CarParkLotRepository interface {
	GetAll() ([]CarParkLot, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]CarParkLot, error)
	Count() (int, error)
	GetByID(id int64) (*CarParkLot, error)
	Create(l *CarParkLot) error
	Update(l *CarParkLot) error
	Delete(id int64) error
}

type mysqlCarParkLotRepository struct {
	db *sql.DB
}

func NewCarParkLotRepository(db *sql.DB) CarParkLotRepository {
	return &mysqlCarParkLotRepository{db: db}
}

const carParkLotColumns = `
	cpl.id, cpl.car_park_id,
	COALESCE(cp.name, '') AS car_park_name,
	cpl.lot_number, cpl.level, cpl.status,
	cpl.created_by, cpl.created_at, cpl.updated_by, cpl.updated_at`

const carParkLotFrom = `
	FROM car_park_lot cpl
	LEFT JOIN car_park cp ON cpl.car_park_id = cp.id`

func scanCarParkLot(row interface {
	Scan(dest ...any) error
}) (*CarParkLot, error) {
	var l CarParkLot
	err := row.Scan(
		&l.ID, &l.CarParkID, &l.CarParkName,
		&l.LotNumber, &l.Level, &l.Status,
		&l.CreatedBy, &l.CreatedAt, &l.UpdatedBy, &l.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *mysqlCarParkLotRepository) GetAll() ([]CarParkLot, error) {
	query := "SELECT" + carParkLotColumns + carParkLotFrom + " ORDER BY cpl.id DESC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lots []CarParkLot
	for rows.Next() {
		l, err := scanCarParkLot(rows)
		if err != nil {
			return nil, err
		}
		lots = append(lots, *l)
	}
	return lots, nil
}

func (r *mysqlCarParkLotRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]CarParkLot, error) {
	sortableColumns := map[string]bool{
		"id":         true,
		"lot_number": true,
		"level":      true,
		"status":     true,
		"updated_by": true,
		"updated_at": true,
	}
	if !sortableColumns[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := "SELECT" + carParkLotColumns + carParkLotFrom +
		" ORDER BY cpl." + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lots []CarParkLot
	for rows.Next() {
		l, err := scanCarParkLot(rows)
		if err != nil {
			return nil, err
		}
		lots = append(lots, *l)
	}
	return lots, nil
}

func (r *mysqlCarParkLotRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM car_park_lot").Scan(&count)
	return count, err
}

func (r *mysqlCarParkLotRepository) GetByID(id int64) (*CarParkLot, error) {
	query := "SELECT" + carParkLotColumns + carParkLotFrom + " WHERE cpl.id = ?"
	row := r.db.QueryRow(query, id)
	l, err := scanCarParkLot(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return l, err
}

func (r *mysqlCarParkLotRepository) Create(l *CarParkLot) error {
	res, err := r.db.Exec(
		`INSERT INTO car_park_lot
			(car_park_id, lot_number, level, status, created_by, created_at, updated_by, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		l.CarParkID, l.LotNumber, l.Level, l.Status,
		l.CreatedBy, l.CreatedAt, l.UpdatedBy, l.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	l.ID = id
	return nil
}

func (r *mysqlCarParkLotRepository) Update(l *CarParkLot) error {
	_, err := r.db.Exec(
		`UPDATE car_park_lot SET
			car_park_id = ?, lot_number = ?, level = ?, status = ?,
			updated_by = ?, updated_at = ?
		 WHERE id = ?`,
		l.CarParkID, l.LotNumber, l.Level, l.Status,
		l.UpdatedBy, l.UpdatedAt, l.ID,
	)
	return err
}

func (r *mysqlCarParkLotRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM car_park_lot WHERE id = ?", id)
	return err
}
