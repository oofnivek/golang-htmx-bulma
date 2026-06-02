package vehicle

import (
	"database/sql"
	"fmt"
)

type FuelCardRepository interface {
	GetAll() ([]FuelCard, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]FuelCard, error)
	Count() (int, error)
	GetByID(id int64) (*FuelCard, error)
	Create(fc *FuelCard) error
	Update(fc *FuelCard) error
	Delete(id int64) error
}

type mysqlFuelCardRepository struct {
	db *sql.DB
}

func NewFuelCardRepository(db *sql.DB) FuelCardRepository {
	return &mysqlFuelCardRepository{db: db}
}

const fuelCardColumns = `
	fc.id, fc.card_no,
	fc.fuel_company_id, COALESCE(fco.name, '') AS fuel_company_name,
	fc.pin_number,
	fc.vehicle_id, v.plate_number,
	fc.status,
	fc.created_by, fc.created_at, fc.updated_by, fc.updated_at`

const fuelCardFrom = `
	FROM fuel_card fc
	LEFT JOIN fuel_company fco ON fc.fuel_company_id = fco.id
	LEFT JOIN vehicle v ON fc.vehicle_id = v.id`

func scanFuelCard(row interface{ Scan(dest ...any) error }) (*FuelCard, error) {
	var fc FuelCard
	err := row.Scan(
		&fc.ID, &fc.CardNo,
		&fc.FuelCompanyID, &fc.FuelCompanyName,
		&fc.PinNumber,
		&fc.VehicleID, &fc.PlateNumber,
		&fc.Status,
		&fc.CreatedBy, &fc.CreatedAt, &fc.UpdatedBy, &fc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &fc, nil
}

var fuelCardSortableColumns = map[string]bool{
	"id":         true,
	"card_no":    true,
	"status":     true,
	"updated_by": true,
	"updated_at": true,
}

func (r *mysqlFuelCardRepository) GetAll() ([]FuelCard, error) {
	query := "SELECT" + fuelCardColumns + fuelCardFrom + " ORDER BY fc.id DESC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []FuelCard
	for rows.Next() {
		fc, err := scanFuelCard(rows)
		if err != nil {
			return nil, err
		}
		cards = append(cards, *fc)
	}
	return cards, nil
}

func (r *mysqlFuelCardRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]FuelCard, error) {
	if !fuelCardSortableColumns[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := fmt.Sprintf(
		"SELECT%s%s ORDER BY fc.%s %s LIMIT ? OFFSET ?",
		fuelCardColumns, fuelCardFrom, sortBy, sortOrder,
	)
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []FuelCard
	for rows.Next() {
		fc, err := scanFuelCard(rows)
		if err != nil {
			return nil, err
		}
		cards = append(cards, *fc)
	}
	return cards, nil
}

func (r *mysqlFuelCardRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM fuel_card").Scan(&count)
	return count, err
}

func (r *mysqlFuelCardRepository) GetByID(id int64) (*FuelCard, error) {
	query := "SELECT" + fuelCardColumns + fuelCardFrom + " WHERE fc.id = ?"
	row := r.db.QueryRow(query, id)
	fc, err := scanFuelCard(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return fc, err
}

func (r *mysqlFuelCardRepository) Create(fc *FuelCard) error {
	res, err := r.db.Exec(
		`INSERT INTO fuel_card
			(card_no, fuel_company_id, pin_number, vehicle_id, status, created_by, created_at, updated_by, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		fc.CardNo, fc.FuelCompanyID, fc.PinNumber, fc.VehicleID, fc.Status,
		fc.CreatedBy, fc.CreatedAt, fc.UpdatedBy, fc.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	fc.ID = id
	return nil
}

func (r *mysqlFuelCardRepository) Update(fc *FuelCard) error {
	_, err := r.db.Exec(
		`UPDATE fuel_card SET
			card_no = ?, fuel_company_id = ?, pin_number = ?, vehicle_id = ?, status = ?,
			updated_by = ?, updated_at = ?
		 WHERE id = ?`,
		fc.CardNo, fc.FuelCompanyID, fc.PinNumber, fc.VehicleID, fc.Status,
		fc.UpdatedBy, fc.UpdatedAt, fc.ID,
	)
	return err
}

func (r *mysqlFuelCardRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM fuel_card WHERE id = ?", id)
	return err
}
