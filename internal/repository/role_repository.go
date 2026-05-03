package repository

import (
	"database/sql"
	"golang-htmx-bulma/internal/model"
)

type RoleRepository interface {
	GetAll() ([]model.Role, error)
	GetPaged(limit, offset int, sortBy, sortOrder, search string) ([]model.Role, error)
	Count(search string) (int, error)
	GetByID(id string) (*model.Role, error)
	Create(role *model.Role) error
	Update(role *model.Role) error
	Delete(id string) error
}

type mysqlRoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &mysqlRoleRepository{db: db}
}

func (r *mysqlRoleRepository) GetAll() ([]model.Role, error) {
	rows, err := r.db.Query("SELECT id, name FROM roles ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *mysqlRoleRepository) GetPaged(limit, offset int, sortBy, sortOrder, search string) ([]model.Role, error) {
	validColumns := map[string]bool{
		"id":   true,
		"name": true,
	}
	if !validColumns[sortBy] {
		sortBy = "id"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	query := "SELECT id, name FROM roles"
	var args []interface{}

	if len(search) >= 2 {
		// Note: The user didn't specify if there's a FULLTEXT index on 'name'.
		// AGENTS.md says to avoid LIKE '%keyword%' and use MATCH() AGAINST().
		// I will assume there is or should be a FULLTEXT index on 'name'.
		query += " WHERE MATCH(name) AGAINST(? IN BOOLEAN MODE)"
		args = append(args, search)
	}

	query += " ORDER BY " + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *mysqlRoleRepository) Count(search string) (int, error) {
	query := "SELECT COUNT(*) FROM roles"
	var args []interface{}
	if len(search) >= 2 {
		query += " WHERE MATCH(name) AGAINST(? IN BOOLEAN MODE)"
		args = append(args, search)
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

func (r *mysqlRoleRepository) GetByID(id string) (*model.Role, error) {
	row := r.db.QueryRow("SELECT id, name FROM roles WHERE id = ?", id)
	var role model.Role
	err := row.Scan(&role.ID, &role.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *mysqlRoleRepository) Create(role *model.Role) error {
	_, err := r.db.Exec("INSERT INTO roles (id, name) VALUES (?, ?)", role.ID, role.Name)
	return err
}

func (r *mysqlRoleRepository) Update(role *model.Role) error {
	_, err := r.db.Exec("UPDATE roles SET name = ? WHERE id = ?", role.Name, role.ID)
	return err
}

func (r *mysqlRoleRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM roles WHERE id = ?", id)
	return err
}
