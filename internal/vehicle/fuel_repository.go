package vehicle

import (
	"database/sql"
)

type VehicleFuelRepository interface {
	GetAll() ([]VehicleFuel, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]VehicleFuel, error)
	Count() (int, error)
	GetByID(id int64) (*VehicleFuel, error)
	Create(f *VehicleFuel) error
	Update(f *VehicleFuel) error
	Delete(id int64) error
}

type mysqlVehicleFuelRepository struct {
	db *sql.DB
}

func NewVehicleFuelRepository(db *sql.DB) VehicleFuelRepository {
	return &mysqlVehicleFuelRepository{db: db}
}

const fuelSelectCols = `
	vf.id, vf.status, vf.vehicle_make_id, vf.vehicle_model_id, vf.fuel_type_id,
	vmk.name, vmdl.name, ft.name,
	vf.fuel_tank_size, vf.fuel_consumption,
	vf.created_by, vf.created_at,
	vf.updated_by, vf.updated_at`

const fuelFromJoin = `
	FROM vehicle_fuel vf
	INNER JOIN vehicle_make vmk ON vmk.id = vf.vehicle_make_id
	INNER JOIN vehicle_model vmdl ON vmdl.id = vf.vehicle_model_id
	INNER JOIN fuel_type ft ON ft.id = vf.fuel_type_id`

func scanFuel(s interface {
	Scan(...any) error
}) (VehicleFuel, error) {
	var f VehicleFuel
	err := s.Scan(
		&f.ID, &f.Status, &f.VehicleMakeID, &f.VehicleModelID, &f.FuelTypeID,
		&f.VehicleMakeName, &f.VehicleModelName, &f.FuelTypeName,
		&f.FuelTankSize, &f.FuelConsumption,
		&f.CreatedBy, &f.CreatedAt,
		&f.UpdatedBy, &f.UpdatedAt,
	)
	return f, err
}

func (r *mysqlVehicleFuelRepository) GetAll() ([]VehicleFuel, error) {
	query := "SELECT" + fuelSelectCols + fuelFromJoin + " ORDER BY vf.id DESC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fuels []VehicleFuel
	for rows.Next() {
		f, err := scanFuel(rows)
		if err != nil {
			return nil, err
		}
		fuels = append(fuels, f)
	}
	return fuels, nil
}

func (r *mysqlVehicleFuelRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]VehicleFuel, error) {
	colMap := map[string]string{
		"id":                 "vf.id",
		"vehicle_make_name":  "vmk.name",
		"vehicle_model_name": "vmdl.name",
		"fuel_type_name":     "ft.name",
		"fuel_tank_size":     "vf.fuel_tank_size",
		"fuel_consumption":   "vf.fuel_consumption",
		"status":             "vf.status",
		"updated_by":         "vf.updated_by",
		"updated_at":         "vf.updated_at",
	}
	orderCol, ok := colMap[sortBy]
	if !ok {
		orderCol = "vf.id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := "SELECT" + fuelSelectCols + fuelFromJoin +
		" ORDER BY " + orderCol + " " + sortOrder + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fuels []VehicleFuel
	for rows.Next() {
		f, err := scanFuel(rows)
		if err != nil {
			return nil, err
		}
		fuels = append(fuels, f)
	}
	return fuels, nil
}

func (r *mysqlVehicleFuelRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM vehicle_fuel").Scan(&count)
	return count, err
}

func (r *mysqlVehicleFuelRepository) GetByID(id int64) (*VehicleFuel, error) {
	query := "SELECT" + fuelSelectCols + fuelFromJoin + " WHERE vf.id = ?"
	row := r.db.QueryRow(query, id)
	f, err := scanFuel(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *mysqlVehicleFuelRepository) Create(f *VehicleFuel) error {
	res, err := r.db.Exec(
		`INSERT INTO vehicle_fuel
			(status, vehicle_make_id, vehicle_model_id, fuel_type_id,
			 fuel_tank_size, fuel_consumption,
			 created_by, created_at, updated_by, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		f.Status, f.VehicleMakeID, f.VehicleModelID, f.FuelTypeID,
		f.FuelTankSize, f.FuelConsumption,
		f.CreatedBy, f.CreatedAt, f.UpdatedBy, f.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	f.ID = id
	return nil
}

func (r *mysqlVehicleFuelRepository) Update(f *VehicleFuel) error {
	_, err := r.db.Exec(
		`UPDATE vehicle_fuel
			SET status = ?, vehicle_make_id = ?, vehicle_model_id = ?, fuel_type_id = ?,
			    fuel_tank_size = ?, fuel_consumption = ?,
			    updated_by = ?, updated_at = ?
			WHERE id = ?`,
		f.Status, f.VehicleMakeID, f.VehicleModelID, f.FuelTypeID,
		f.FuelTankSize, f.FuelConsumption,
		f.UpdatedBy, f.UpdatedAt, f.ID,
	)
	return err
}

func (r *mysqlVehicleFuelRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM vehicle_fuel WHERE id = ?", id)
	return err
}
