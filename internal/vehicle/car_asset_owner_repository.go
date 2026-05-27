package vehicle

import (
	"database/sql"
)

type CarAssetOwnerRepository interface {
	GetAll() ([]CarAssetOwner, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]CarAssetOwner, error)
	Count() (int, error)
	GetByID(id int64) (*CarAssetOwner, error)
	Create(owner *CarAssetOwner) error
	Update(owner *CarAssetOwner) error
	Delete(id int64) error
}

type mysqlCarAssetOwnerRepository struct {
	db *sql.DB
}

func NewCarAssetOwnerRepository(db *sql.DB) CarAssetOwnerRepository {
	return &mysqlCarAssetOwnerRepository{db: db}
}

func (r *mysqlCarAssetOwnerRepository) GetAll() ([]CarAssetOwner, error) {
	rows, err := r.db.Query("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM car_asset_owner ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var owners []CarAssetOwner
	for rows.Next() {
		var o CarAssetOwner
		if err := rows.Scan(&o.ID, &o.Name, &o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedBy, &o.UpdatedAt); err != nil {
			return nil, err
		}
		owners = append(owners, o)
	}
	return owners, nil
}

func (r *mysqlCarAssetOwnerRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]CarAssetOwner, error) {
	sortableColumns := map[string]bool{
		"id":         true,
		"name":       true,
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

	query := "SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM car_asset_owner ORDER BY " + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var owners []CarAssetOwner
	for rows.Next() {
		var o CarAssetOwner
		if err := rows.Scan(&o.ID, &o.Name, &o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedBy, &o.UpdatedAt); err != nil {
			return nil, err
		}
		owners = append(owners, o)
	}
	return owners, nil
}

func (r *mysqlCarAssetOwnerRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM car_asset_owner").Scan(&count)
	return count, err
}

func (r *mysqlCarAssetOwnerRepository) GetByID(id int64) (*CarAssetOwner, error) {
	row := r.db.QueryRow("SELECT id, name, status, created_by, created_at, updated_by, updated_at FROM car_asset_owner WHERE id = ?", id)
	var o CarAssetOwner
	err := row.Scan(&o.ID, &o.Name, &o.Status, &o.CreatedBy, &o.CreatedAt, &o.UpdatedBy, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *mysqlCarAssetOwnerRepository) Create(o *CarAssetOwner) error {
	res, err := r.db.Exec("INSERT INTO car_asset_owner (name, status, created_by, created_at, updated_by, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		o.Name, o.Status, o.CreatedBy, o.CreatedAt, o.UpdatedBy, o.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	o.ID = id
	return nil
}

func (r *mysqlCarAssetOwnerRepository) Update(o *CarAssetOwner) error {
	_, err := r.db.Exec("UPDATE car_asset_owner SET name = ?, status = ?, updated_by = ?, updated_at = ? WHERE id = ?",
		o.Name, o.Status, o.UpdatedBy, o.UpdatedAt, o.ID)
	return err
}

func (r *mysqlCarAssetOwnerRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM car_asset_owner WHERE id = ?", id)
	return err
}
