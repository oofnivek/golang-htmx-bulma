package user

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// User represents the user model in the system.
type User struct {
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Mobile       string    `json:"mobile"`
	Designation  string    `json:"designation"`
	Department   string    `json:"department"`
	IsEnabled    bool      `json:"is_enabled"`
	CreatedAt    time.Time `json:"created_at_utc"`
	UpdatedAt    time.Time `json:"updated_at_utc"`
	RoleID       string    `json:"role_id"`

	PasswordHash *string `json:"-"`
	Password     string `json:"-"`
	ConfirmPass  string `json:"-"`

	// Joining with roles table for display
	RoleName string `json:"role_name"`
}

// Role represents a user role in the system.
type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CustomClaims extends standard JWT claims with app-specific info
type CustomClaims struct {
	Email     string `json:"email"`
	RoleID    string `json:"role_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	jwt.RegisteredClaims
}
