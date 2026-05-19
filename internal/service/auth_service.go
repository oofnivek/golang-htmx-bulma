package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"golang-htmx-bulma/internal/model"
	"golang-htmx-bulma/internal/pkg/crypto"
)

// CustomClaims extends standard JWT claims with app-specific info
type CustomClaims struct {
	Email     string `json:"email"`
	RoleID    string `json:"role_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Login(email, password string) (string, *model.User, error)
}

type authService struct {
	userSvc    UserService
	signingKey []byte
}

func NewAuthService(userSvc UserService, signingKey string) AuthService {
	return &authService{
		userSvc:    userSvc,
		signingKey: []byte(signingKey),
	}
}

func (s *authService) Login(email, password string) (string, *model.User, error) {
	// 1. Get user by email
	user, err := s.userSvc.GetByEmail(email)
	if err != nil {
		// Log the actual internal database/system error securely on the server
		slog.Error("Database query failed during login", "email", email, "error", err.Error())
		return "", nil, errors.New("invalid email or password")
	}
	if user == nil {
		slog.Warn("Login failed: email not found", "email", email)
		return "", nil, errors.New("invalid email or password")
	}

	// 2. Check if user is enabled
	if !user.IsEnabled {
		slog.Warn("Login failed: account is disabled", "email", email)
		return "", nil, errors.New("account is disabled")
	}

	// 3. Verify password
	if err := crypto.VerifyPasswordV3(password, user.PasswordHash); err != nil {
		slog.Warn("Login failed: incorrect password", "email", email)
		return "", nil, errors.New("invalid email or password")
	}

	// 4. Generate JWT ID
	jti, err := generateJTI()
	if err != nil {
		slog.Error("Failed to generate token ID during login", "email", email, "error", err.Error())
		return "", nil, fmt.Errorf("failed to generate token id: %w", err)
	}

	// 5. Create JWT claims
	now := time.Now().UTC()
	claims := CustomClaims{
		Email:     user.Email,
		RoleID:    user.RoleID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email,
			Issuer:    "fleet-management-system",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)), // 24 hours expiry
			ID:        jti,
		},
	}

	// 6. Sign JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.signingKey)
	if err != nil {
		slog.Error("Failed to sign token during login", "email", email, "error", err.Error())
		return "", nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, user, nil
}

// generateJTI creates a random UUID v4 string to uniquely identify the token.
func generateJTI() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Set version 4 and variant bits per RFC 4122.
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
