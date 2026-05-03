package model

import "time"

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
	
	// Password related fields ignored for now
	// PasswordResetRequestID []byte
	// PasswordHash           string

	// Joining with roles table for display
	RoleName string `json:"role_name"`
}
