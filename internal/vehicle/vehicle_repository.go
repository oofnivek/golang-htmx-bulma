package vehicle

import (
	"database/sql"
)

type VehicleRepository interface {
	GetAll() ([]Vehicle, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]Vehicle, error)
	Count() (int, error)
	GetByID(id int64) (*Vehicle, error)
	Create(v *Vehicle) error
	Update(v *Vehicle) error
	Delete(id int64) error
}

type mysqlVehicleRepository struct {
	db *sql.DB
}

func NewVehicleRepository(db *sql.DB) VehicleRepository {
	return &mysqlVehicleRepository{db: db}
}

const vehicleColumns = `
	v.id,
	v.vehicle_make_id, COALESCE(vm.name, '') AS vehicle_make_name,
	v.vehicle_model_id, COALESCE(vmd.name, '') AS vehicle_model_name,
	v.vehicle_type_id, COALESCE(vt.name, '') AS vehicle_type_name,
	v.fuel_type_id, COALESCE(ft.name, '') AS fuel_type_name,
	v.vehicle_color_id, COALESCE(vc.name, '') AS vehicle_color_name,
	v.description, v.plate_number, v.iu_number, v.chassis_number, v.engine_number,
	v.num_seats, v.boot_space,
	v.car_park_id, COALESCE(cp.name, '') AS car_park_name,
	v.asset_owner_id, COALESCE(cao.name, '') AS asset_owner_name,
	v.last_service_date, v.last_cleaned_date,
	v.last_service_mileage, v.current_mileage, v.current_fuel_level,
	v.active_from, v.active_to,
	v.vehicle_status_id, COALESCE(vs.substatus, '') AS vehicle_status_name,
	v.created_by, v.created_at, v.updated_by, v.updated_at`

const vehicleFrom = `
	FROM vehicle v
	LEFT JOIN vehicle_make vm ON v.vehicle_make_id = vm.id
	LEFT JOIN vehicle_model vmd ON v.vehicle_model_id = vmd.id
	LEFT JOIN vehicle_type vt ON v.vehicle_type_id = vt.id
	LEFT JOIN fuel_type ft ON v.fuel_type_id = ft.id
	LEFT JOIN vehicle_color vc ON v.vehicle_color_id = vc.id
	LEFT JOIN car_park cp ON v.car_park_id = cp.id
	LEFT JOIN car_asset_owner cao ON v.asset_owner_id = cao.id
	LEFT JOIN vehicle_status vs ON v.vehicle_status_id = vs.id`

func scanVehicle(row interface {
	Scan(dest ...any) error
}) (*Vehicle, error) {
	var v Vehicle
	err := row.Scan(
		&v.ID,
		&v.VehicleMakeID, &v.VehicleMakeName,
		&v.VehicleModelID, &v.VehicleModelName,
		&v.VehicleTypeID, &v.VehicleTypeName,
		&v.FuelTypeID, &v.FuelTypeName,
		&v.VehicleColorID, &v.VehicleColorName,
		&v.Description, &v.PlateNumber, &v.IUNumber, &v.ChassisNumber, &v.EngineNumber,
		&v.NumSeats, &v.BootSpace,
		&v.CarParkID, &v.CarParkName,
		&v.AssetOwnerID, &v.AssetOwnerName,
		&v.LastServiceDate, &v.LastCleanedDate,
		&v.LastServiceMileage, &v.CurrentMileage, &v.CurrentFuelLevel,
		&v.ActiveFrom, &v.ActiveTo,
		&v.VehicleStatusID, &v.VehicleStatusName,
		&v.CreatedBy, &v.CreatedAt, &v.UpdatedBy, &v.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *mysqlVehicleRepository) GetAll() ([]Vehicle, error) {
	query := "SELECT" + vehicleColumns + vehicleFrom + " ORDER BY v.id DESC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		v, err := scanVehicle(rows)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, *v)
	}
	return vehicles, nil
}

func (r *mysqlVehicleRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]Vehicle, error) {
	sortableColumns := map[string]bool{
		"id":                true,
		"plate_number":      true,
		"vehicle_make_id":   true,
		"vehicle_model_id":  true,
		"vehicle_status_id": true,
		"updated_by":        true,
		"updated_at":        true,
	}
	if !sortableColumns[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := "SELECT" + vehicleColumns + vehicleFrom +
		" ORDER BY v." + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		v, err := scanVehicle(rows)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, *v)
	}
	return vehicles, nil
}

func (r *mysqlVehicleRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM vehicle").Scan(&count)
	return count, err
}

func (r *mysqlVehicleRepository) GetByID(id int64) (*Vehicle, error) {
	query := "SELECT" + vehicleColumns + vehicleFrom + " WHERE v.id = ?"
	row := r.db.QueryRow(query, id)
	v, err := scanVehicle(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return v, err
}

func (r *mysqlVehicleRepository) Create(v *Vehicle) error {
	res, err := r.db.Exec(
		`INSERT INTO vehicle
			(vehicle_make_id, vehicle_model_id, vehicle_type_id, fuel_type_id, vehicle_color_id,
			 description, plate_number, iu_number, chassis_number, engine_number,
			 num_seats, boot_space, car_park_id, asset_owner_id, vehicle_status_id,
			 last_service_date, last_cleaned_date, last_service_mileage, current_mileage, current_fuel_level,
			 active_from, active_to, created_by, created_at, updated_by, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		v.VehicleMakeID, v.VehicleModelID, v.VehicleTypeID, v.FuelTypeID, v.VehicleColorID,
		v.Description, v.PlateNumber, v.IUNumber, v.ChassisNumber, v.EngineNumber,
		v.NumSeats, v.BootSpace, v.CarParkID, v.AssetOwnerID, v.VehicleStatusID,
		v.LastServiceDate, v.LastCleanedDate, v.LastServiceMileage, v.CurrentMileage, v.CurrentFuelLevel,
		v.ActiveFrom, v.ActiveTo, v.CreatedBy, v.CreatedAt, v.UpdatedBy, v.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	v.ID = id
	return nil
}

func (r *mysqlVehicleRepository) Update(v *Vehicle) error {
	_, err := r.db.Exec(
		`UPDATE vehicle SET
			vehicle_make_id = ?, vehicle_model_id = ?, vehicle_type_id = ?, fuel_type_id = ?, vehicle_color_id = ?,
			description = ?, plate_number = ?, iu_number = ?, chassis_number = ?, engine_number = ?,
			num_seats = ?, boot_space = ?, car_park_id = ?, asset_owner_id = ?, vehicle_status_id = ?,
			last_service_date = ?, last_cleaned_date = ?, last_service_mileage = ?, current_mileage = ?, current_fuel_level = ?,
			active_from = ?, active_to = ?, updated_by = ?, updated_at = ?
		 WHERE id = ?`,
		v.VehicleMakeID, v.VehicleModelID, v.VehicleTypeID, v.FuelTypeID, v.VehicleColorID,
		v.Description, v.PlateNumber, v.IUNumber, v.ChassisNumber, v.EngineNumber,
		v.NumSeats, v.BootSpace, v.CarParkID, v.AssetOwnerID, v.VehicleStatusID,
		v.LastServiceDate, v.LastCleanedDate, v.LastServiceMileage, v.CurrentMileage, v.CurrentFuelLevel,
		v.ActiveFrom, v.ActiveTo, v.UpdatedBy, v.UpdatedAt, v.ID,
	)
	return err
}

func (r *mysqlVehicleRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM vehicle WHERE id = ?", id)
	return err
}
