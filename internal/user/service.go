package user

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"golang-htmx-bulma/internal/pkg/crypto"
)

// UserService defines business operations on User.
type UserService interface {
	ListAll() ([]User, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]User, int, error)
	GetByEmail(email string) (*User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(email string) error
}

type localUserService struct {
	repo UserRepository
}

// NewLocalUserService creates a new in-memory/DB-backed UserService.
func NewLocalUserService(repo UserRepository) UserService {
	return &localUserService{repo: repo}
}

func (s *localUserService) ListAll() ([]User, error) {
	return s.repo.GetAll()
}

func (s *localUserService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]User, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	users, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *localUserService) GetByEmail(email string) (*User, error) {
	return s.repo.GetByEmail(email)
}

func (s *localUserService) CreateUser(u *User) error {
	if u.Password != u.ConfirmPass {
		return fmt.Errorf("passwords do not match")
	}

	hash, err := crypto.HashPasswordV3(u.Password, 100000)
	if err != nil {
		return err
	}
	u.PasswordHash = &hash

	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now
	return s.repo.Create(u)
}

func (s *localUserService) UpdateUser(u *User) error {
	if u.Password != "" {
		if u.Password != u.ConfirmPass {
			return fmt.Errorf("passwords do not match")
		}

		hash, err := crypto.HashPasswordV3(u.Password, 100000)
		if err != nil {
			return err
		}
		u.PasswordHash = &hash
	}

	u.UpdatedAt = time.Now().UTC()
	return s.repo.Update(u)
}

func (s *localUserService) DeleteUser(email string) error {
	return s.repo.Delete(email)
}

// RoleService defines business operations on Role.
type RoleService interface {
	ListAll() ([]Role, error)
	ListPaged(page, pageSize int, sortBy, sortOrder string) ([]Role, int, error)
	FindByID(id string) (*Role, error)
	CreateRole(id, name string) (*Role, error)
	UpdateRole(id, name string) (*Role, error)
	DeleteRole(id string) error
}

type localRoleService struct {
	repo RoleRepository
}

// NewLocalRoleService creates a new in-memory/DB-backed RoleService.
func NewLocalRoleService(repo RoleRepository) RoleService {
	return &localRoleService{repo: repo}
}

func (s *localRoleService) ListAll() ([]Role, error) {
	return s.repo.GetAll()
}

func (s *localRoleService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]Role, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	roles, err := s.repo.GetPaged(pageSize, offset, sortBy, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (s *localRoleService) FindByID(id string) (*Role, error) {
	return s.repo.GetByID(id)
}

func (s *localRoleService) CreateRole(id, name string) (*Role, error) {
	role := &Role{
		ID:   id,
		Name: name,
	}
	err := s.repo.Create(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *localRoleService) UpdateRole(id, name string) (*Role, error) {
	role, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	role.Name = name

	err = s.repo.Update(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *localRoleService) DeleteRole(id string) error {
	return s.repo.Delete(id)
}

// AuthService defines business operations for user authentication and session issuing.
type AuthService interface {
	Login(email, password string) (string, *User, error)
}

type localAuthService struct {
	userSvc    UserService
	signingKey []byte
}

// NewLocalAuthService creates a new in-memory/DB-backed AuthService.
func NewLocalAuthService(userSvc UserService, signingKey string) AuthService {
	return &localAuthService{
		userSvc:    userSvc,
		signingKey: []byte(signingKey),
	}
}

func (s *localAuthService) Login(email, password string) (string, *User, error) {
	user, err := s.userSvc.GetByEmail(email)
	if err != nil {
		slog.Error("Database query failed during login", "email", email, "error", err.Error())
		return "", nil, errors.New("invalid email or password")
	}
	if user == nil {
		slog.Warn("Login failed: email not found", "email", email)
		return "", nil, errors.New("invalid email or password")
	}

	if !user.IsEnabled {
		slog.Warn("Login failed: account is disabled", "email", email)
		return "", nil, errors.New("account is disabled")
	}

	passwordHash := ""
	if user.PasswordHash != nil {
		passwordHash = *user.PasswordHash
	}
	verificationResult := crypto.VerifyUserPassword(email, passwordHash, password)
	if verificationResult == crypto.PasswordVerificationFailed {
		slog.Warn("Login failed: incorrect password", "email", email)
		return "", nil, errors.New("invalid email or password")
	}

	if verificationResult == crypto.PasswordVerificationSuccessRehash {
		slog.Info("Password verification succeeded using legacy fallback; rehash recommended", "email", email)
	}

	jti, err := generateJTI()
	if err != nil {
		slog.Error("Failed to generate token ID during login", "email", email, "error", err.Error())
		return "", nil, fmt.Errorf("failed to generate token id: %w", err)
	}

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
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.signingKey)
	if err != nil {
		slog.Error("Failed to sign token during login", "email", email, "error", err.Error())
		return "", nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, user, nil
}

func generateJTI() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
