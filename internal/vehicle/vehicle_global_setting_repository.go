package vehicle

import (
	"database/sql"
	"fmt"
)

type VehicleGlobalSettingRepository interface {
	GetAll() ([]VehicleGlobalSetting, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]VehicleGlobalSetting, error)
	Count() (int, error)
	GetByID(id int64) (*VehicleGlobalSetting, error)
	Create(s *VehicleGlobalSetting) error
	Update(s *VehicleGlobalSetting) error
	Delete(id int64) error
}

type mysqlVehicleGlobalSettingRepository struct {
	db *sql.DB
}

func NewVehicleGlobalSettingRepository(db *sql.DB) VehicleGlobalSettingRepository {
	return &mysqlVehicleGlobalSettingRepository{db: db}
}

const vehicleGlobalSettingColumns = "id, `key`, value, remark, country_code, created_by, created_at, updated_by, updated_at"
const vehicleGlobalSettingFrom = "FROM global_setting"

func scanVehicleGlobalSetting(s interface{ Scan(...any) error }) (*VehicleGlobalSetting, error) {
	var g VehicleGlobalSetting
	err := s.Scan(&g.ID, &g.Key, &g.Value, &g.Remark, &g.CountryCode, &g.CreatedBy, &g.CreatedAt, &g.UpdatedBy, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

var vehicleGlobalSettingSortableColumns = map[string]bool{
	"id":           true,
	"key":          true,
	"value":        true,
	"country_code": true,
	"created_at":   true,
	"updated_at":   true,
}

func (r *mysqlVehicleGlobalSettingRepository) GetAll() ([]VehicleGlobalSetting, error) {
	rows, err := r.db.Query("SELECT " + vehicleGlobalSettingColumns + " " + vehicleGlobalSettingFrom + " ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []VehicleGlobalSetting
	for rows.Next() {
		g, err := scanVehicleGlobalSetting(rows)
		if err != nil {
			return nil, err
		}
		settings = append(settings, *g)
	}
	return settings, nil
}

func (r *mysqlVehicleGlobalSettingRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]VehicleGlobalSetting, error) {
	if !vehicleGlobalSettingSortableColumns[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := fmt.Sprintf("SELECT %s %s ORDER BY %s %s LIMIT ? OFFSET ?", vehicleGlobalSettingColumns, vehicleGlobalSettingFrom, sortBy, sortOrder)
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []VehicleGlobalSetting
	for rows.Next() {
		g, err := scanVehicleGlobalSetting(rows)
		if err != nil {
			return nil, err
		}
		settings = append(settings, *g)
	}
	return settings, nil
}

func (r *mysqlVehicleGlobalSettingRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) " + vehicleGlobalSettingFrom).Scan(&count)
	return count, err
}

func (r *mysqlVehicleGlobalSettingRepository) GetByID(id int64) (*VehicleGlobalSetting, error) {
	row := r.db.QueryRow("SELECT "+vehicleGlobalSettingColumns+" "+vehicleGlobalSettingFrom+" WHERE id = ?", id)
	g, err := scanVehicleGlobalSetting(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return g, err
}

func (r *mysqlVehicleGlobalSettingRepository) Create(s *VehicleGlobalSetting) error {
	res, err := r.db.Exec(
		"INSERT INTO global_setting (`key`, value, remark, country_code, created_by, created_at, updated_by, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		s.Key, s.Value, s.Remark, s.CountryCode, s.CreatedBy, s.CreatedAt, s.UpdatedBy, s.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	s.ID = id
	return nil
}

func (r *mysqlVehicleGlobalSettingRepository) Update(s *VehicleGlobalSetting) error {
	_, err := r.db.Exec(
		"UPDATE global_setting SET `key` = ?, value = ?, remark = ?, country_code = ?, updated_by = ?, updated_at = ? WHERE id = ?",
		s.Key, s.Value, s.Remark, s.CountryCode, s.UpdatedBy, s.UpdatedAt, s.ID,
	)
	return err
}

func (r *mysqlVehicleGlobalSettingRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM global_setting WHERE id = ?", id)
	return err
}
