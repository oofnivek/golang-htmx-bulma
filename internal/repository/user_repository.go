package repository

import (
	"database/sql"
	"golang-htmx-bulma/internal/model"
)

type UserRepository interface {
	GetAll() ([]model.User, error)
	GetPaged(limit, offset int, sortBy, sortOrder, search string) ([]model.User, error)
	Count(search string) (int, error)
	GetByEmail(email string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(email string) error
}

type mysqlUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) GetAll() ([]model.User, error) {
	query := `
		SELECT u.email, u.first_name, u.last_name, u.mobile, u.designation, u.department, u.is_enabled, 
		       u.created_at_utc, u.updated_at_utc, u.role_id, u.password_hash, r.name as role_name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		ORDER BY u.email ASC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.Email, &u.FirstName, &u.LastName, &u.Mobile, &u.Designation, &u.Department, &u.IsEnabled,
			&u.CreatedAt, &u.UpdatedAt, &u.RoleID, &u.PasswordHash, &u.RoleName)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *mysqlUserRepository) GetPaged(limit, offset int, sortBy, sortOrder, search string) ([]model.User, error) {
	validColumns := map[string]bool{
		"email":          true,
		"first_name":     true,
		"last_name":      true,
		"designation":    true,
		"department":     true,
		"is_enabled":     true,
		"updated_at_utc": true,
		"role_name":      true,
	}
	
	// Map role_name to r.name for sorting
	sortCol := "u." + sortBy
	if sortBy == "role_name" {
		sortCol = "r.name"
	} else if !validColumns[sortBy] {
		sortCol = "u.email"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	query := `
		SELECT u.email, u.first_name, u.last_name, u.mobile, u.designation, u.department, u.is_enabled, 
		       u.created_at_utc, u.updated_at_utc, u.role_id, u.password_hash, r.name as role_name
		FROM users u
		JOIN roles r ON u.role_id = r.id`
	
	var args []interface{}
	if len(search) >= 2 {
		query += " WHERE MATCH(u.email, u.first_name, u.last_name, u.mobile) AGAINST(? IN BOOLEAN MODE)"
		args = append(args, search)
	}

	query += " ORDER BY " + sortCol + " " + sortOrder + " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.Email, &u.FirstName, &u.LastName, &u.Mobile, &u.Designation, &u.Department, &u.IsEnabled,
			&u.CreatedAt, &u.UpdatedAt, &u.RoleID, &u.PasswordHash, &u.RoleName)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *mysqlUserRepository) Count(search string) (int, error) {
	query := "SELECT COUNT(*) FROM users u"
	var args []interface{}
	if len(search) >= 2 {
		query += " WHERE MATCH(u.email, u.first_name, u.last_name, u.mobile) AGAINST(? IN BOOLEAN MODE)"
		args = append(args, search)
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

func (r *mysqlUserRepository) GetByEmail(email string) (*model.User, error) {
	query := `
		SELECT u.email, u.first_name, u.last_name, u.mobile, u.designation, u.department, u.is_enabled, 
		       u.created_at_utc, u.updated_at_utc, u.role_id, u.password_hash, r.name as role_name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.email = ?`
	
	var u model.User
	err := r.db.QueryRow(query, email).Scan(&u.Email, &u.FirstName, &u.LastName, &u.Mobile, &u.Designation, &u.Department, &u.IsEnabled,
		&u.CreatedAt, &u.UpdatedAt, &u.RoleID, &u.PasswordHash, &u.RoleName)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *mysqlUserRepository) Create(u *model.User) error {
	query := `
		INSERT INTO users (email, first_name, last_name, mobile, designation, department, is_enabled, created_at_utc, updated_at_utc, role_id, password_hash)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := r.db.Exec(query, u.Email, u.FirstName, u.LastName, u.Mobile, u.Designation, u.Department, u.IsEnabled, u.CreatedAt, u.UpdatedAt, u.RoleID, u.PasswordHash)
	return err
}

func (r *mysqlUserRepository) Update(u *model.User) error {
	var err error
	if u.PasswordHash != "" {
		query := `
			UPDATE users SET first_name = ?, last_name = ?, mobile = ?, designation = ?, department = ?, is_enabled = ?, updated_at_utc = ?, role_id = ?, password_hash = ?
			WHERE email = ?`
		_, err = r.db.Exec(query, u.FirstName, u.LastName, u.Mobile, u.Designation, u.Department, u.IsEnabled, u.UpdatedAt, u.RoleID, u.PasswordHash, u.Email)
	} else {
		query := `
			UPDATE users SET first_name = ?, last_name = ?, mobile = ?, designation = ?, department = ?, is_enabled = ?, updated_at_utc = ?, role_id = ?
			WHERE email = ?`
		_, err = r.db.Exec(query, u.FirstName, u.LastName, u.Mobile, u.Designation, u.Department, u.IsEnabled, u.UpdatedAt, u.RoleID, u.Email)
	}
	return err
}

func (r *mysqlUserRepository) Delete(email string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE email = ?", email)
	return err
}
