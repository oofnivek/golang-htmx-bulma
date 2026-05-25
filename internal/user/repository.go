package user

import (
	"database/sql"
)

// UserRepository defines the database operations for User.
type UserRepository interface {
	GetAll() ([]User, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]User, error)
	Count() (int, error)
	GetByEmail(email string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(email string) error
}

type mysqlUserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository instance using MySQL.
func NewUserRepository(db *sql.DB) UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) GetAll() ([]User, error) {
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

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.Email, &u.FirstName, &u.LastName, &u.Mobile, &u.Designation, &u.Department, &u.IsEnabled,
			&u.CreatedAt, &u.UpdatedAt, &u.RoleID, &u.PasswordHash, &u.RoleName)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *mysqlUserRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]User, error) {
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
		JOIN roles r ON u.role_id = r.id
		ORDER BY ` + sortCol + ` ` + sortOrder + ` LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.Email, &u.FirstName, &u.LastName, &u.Mobile, &u.Designation, &u.Department, &u.IsEnabled,
			&u.CreatedAt, &u.UpdatedAt, &u.RoleID, &u.PasswordHash, &u.RoleName)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *mysqlUserRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users u").Scan(&count)
	return count, err
}

func (r *mysqlUserRepository) GetByEmail(email string) (*User, error) {
	query := `
		SELECT u.email, u.first_name, u.last_name, u.mobile, u.designation, u.department, u.is_enabled, 
		       u.created_at_utc, u.updated_at_utc, u.role_id, u.password_hash, r.name as role_name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.email = ?`
	
	var u User
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

func (r *mysqlUserRepository) Create(u *User) error {
	query := `
		INSERT INTO users (email, first_name, last_name, mobile, designation, department, is_enabled, created_at_utc, updated_at_utc, role_id, password_hash)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := r.db.Exec(query, u.Email, u.FirstName, u.LastName, u.Mobile, u.Designation, u.Department, u.IsEnabled, u.CreatedAt, u.UpdatedAt, u.RoleID, u.PasswordHash)
	return err
}

func (r *mysqlUserRepository) Update(u *User) error {
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

// RoleRepository defines the database operations for Role.
type RoleRepository interface {
	GetAll() ([]Role, error)
	GetPaged(limit, offset int, sortBy, sortOrder string) ([]Role, error)
	Count() (int, error)
	GetByID(id string) (*Role, error)
	Create(role *Role) error
	Update(role *Role) error
	Delete(id string) error
}

type mysqlRoleRepository struct {
	db *sql.DB
}

// NewRoleRepository creates a new RoleRepository instance using MySQL.
func NewRoleRepository(db *sql.DB) RoleRepository {
	return &mysqlRoleRepository{db: db}
}

func (r *mysqlRoleRepository) GetAll() ([]Role, error) {
	rows, err := r.db.Query("SELECT id, name FROM roles ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *mysqlRoleRepository) GetPaged(limit, offset int, sortBy, sortOrder string) ([]Role, error) {
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

	query := "SELECT id, name FROM roles ORDER BY " + sortBy + " " + sortOrder + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *mysqlRoleRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM roles").Scan(&count)
	return count, err
}

func (r *mysqlRoleRepository) GetByID(id string) (*Role, error) {
	row := r.db.QueryRow("SELECT id, name FROM roles WHERE id = ?", id)
	var role Role
	err := row.Scan(&role.ID, &role.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *mysqlRoleRepository) Create(role *Role) error {
	_, err := r.db.Exec("INSERT INTO roles (id, name) VALUES (?, ?)", role.ID, role.Name)
	return err
}

func (r *mysqlRoleRepository) Update(role *Role) error {
	_, err := r.db.Exec("UPDATE roles SET name = ? WHERE id = ?", role.Name, role.ID)
	return err
}

func (r *mysqlRoleRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM roles WHERE id = ?", id)
	return err
}
