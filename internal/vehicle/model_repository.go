package vehicle

import (
	"database/sql"
)

type VehicleModelRepository interface {
	GetAll() ([]VehicleModel, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]VehicleModel, error)
	Count() (int, error)
	GetByID(id int64) (*VehicleModel, error)
	Create(m *VehicleModel) error
	Update(m *VehicleModel) error
	Delete(id int64) error
}

type mysqlVehicleModelRepository struct {
	db *sql.DB
}

func NewVehicleModelRepository(db *sql.DB) VehicleModelRepository {
	return &mysqlVehicleModelRepository{db: db}
}

const modelSelectCols = `
	vmod.id, vmod.vehicle_type_id, vmod.vehicle_make_id,
	vt.name, vmk.name,
	vmod.name, vmod.status,
	vmod.created_by, vmod.created_at,
	vmod.updated_by, vmod.updated_at`

const modelFromJoin = `
	FROM vehicle_model vmod
	INNER JOIN vehicle_type vt  ON vt.id  = vmod.vehicle_type_id
	INNER JOIN vehicle_make vmk ON vmk.id = vmod.vehicle_make_id`

func scanModel(s interface {
	Scan(...any) error
}) (VehicleModel, error) {
	var m VehicleModel
	err := s.Scan(
		&m.ID, &m.VehicleTypeID, &m.VehicleMakeID,
		&m.VehicleTypeName, &m.VehicleMakeName,
		&m.Name, &m.Status,
		&m.CreatedBy, &m.CreatedAt,
		&m.UpdatedBy, &m.UpdatedAt,
	)
	return m, err
}

func (r *mysqlVehicleModelRepository) GetAll() ([]VehicleModel, error) {
	query := "SELECT" + modelSelectCols + modelFromJoin + " ORDER BY vmod.id DESC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []VehicleModel
	for rows.Next() {
		m, err := scanModel(rows)
		if err != nil {
			return nil, err
		}
		models = append(models, m)
	}
	return models, nil
}

func (r *mysqlVehicleModelRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]VehicleModel, error) {
	colMap := map[string]string{
		"id":                "vmod.id",
		"name":              "vmod.name",
		"vehicle_type_name": "vt.name",
		"vehicle_make_name": "vmk.name",
		"status":            "vmod.status",
		"updated_by":        "vmod.updated_by",
		"updated_at":        "vmod.updated_at",
	}
	orderCol, ok := colMap[sortBy]
	if !ok {
		orderCol = "vmod.id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := "SELECT" + modelSelectCols + modelFromJoin +
		" ORDER BY " + orderCol + " " + sortOrder + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []VehicleModel
	for rows.Next() {
		m, err := scanModel(rows)
		if err != nil {
			return nil, err
		}
		models = append(models, m)
	}
	return models, nil
}

func (r *mysqlVehicleModelRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM vehicle_model").Scan(&count)
	return count, err
}

func (r *mysqlVehicleModelRepository) GetByID(id int64) (*VehicleModel, error) {
	query := "SELECT" + modelSelectCols + modelFromJoin + " WHERE vmod.id = ?"
	row := r.db.QueryRow(query, id)
	m, err := scanModel(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *mysqlVehicleModelRepository) Create(m *VehicleModel) error {
	res, err := r.db.Exec(
		`INSERT INTO vehicle_model
			(vehicle_type_id, vehicle_make_id, name, status, created_by, created_at, updated_by, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		m.VehicleTypeID, m.VehicleMakeID, m.Name, m.Status,
		m.CreatedBy, m.CreatedAt, m.UpdatedBy, m.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	m.ID = id
	return nil
}

func (r *mysqlVehicleModelRepository) Update(m *VehicleModel) error {
	_, err := r.db.Exec(
		`UPDATE vehicle_model
			SET vehicle_type_id = ?, vehicle_make_id = ?, name = ?, status = ?,
			    updated_by = ?, updated_at = ?
			WHERE id = ?`,
		m.VehicleTypeID, m.VehicleMakeID, m.Name, m.Status,
		m.UpdatedBy, m.UpdatedAt, m.ID,
	)
	return err
}

func (r *mysqlVehicleModelRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM vehicle_model WHERE id = ?", id)
	return err
}
