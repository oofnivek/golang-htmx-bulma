package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type remoteUserService struct {
	client  *http.Client
	baseURL string
}

// NewRemoteUserService creates a new HTTP-based UserService client.
func NewRemoteUserService(baseURL string) UserService {
	return &remoteUserService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteUserService) ListAll() ([]User, error) {
	resp, err := s.client.Get(s.baseURL + "/api/users")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}

	var data struct {
		Users []User `json:"users"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Users, nil
}

func (s *remoteUserService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]User, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)

	u := fmt.Sprintf("%s/api/users?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}

	var data struct {
		Users []User `json:"users"`
		Total int    `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Users, data.Total, nil
}

func (s *remoteUserService) GetByEmail(email string) (*User, error) {
	u := fmt.Sprintf("%s/api/users/%s", s.baseURL, url.PathEscape(email))
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *remoteUserService) CreateUser(u *User) error {
	body, err := json.Marshal(u)
	if err != nil {
		return err
	}

	resp, err := s.client.Post(s.baseURL+"/api/users", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errData map[string]string
		json.NewDecoder(resp.Body).Decode(&errData)
		if msg, ok := errData["error"]; ok {
			return errors.New(msg)
		}
		return fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	return nil
}

func (s *remoteUserService) UpdateUser(u *User) error {
	body, err := json.Marshal(u)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/users/%s", s.baseURL, url.PathEscape(u.Email)), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errData map[string]string
		json.NewDecoder(resp.Body).Decode(&errData)
		if msg, ok := errData["error"]; ok {
			return errors.New(msg)
		}
		return fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	return nil
}

func (s *remoteUserService) DeleteUser(email string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/users/%s", s.baseURL, url.PathEscape(email)), nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errData map[string]string
		json.NewDecoder(resp.Body).Decode(&errData)
		if msg, ok := errData["error"]; ok {
			return errors.New(msg)
		}
		return fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	return nil
}

type remoteRoleService struct {
	client  *http.Client
	baseURL string
}

// NewRemoteRoleService creates a new HTTP-based RoleService client.
func NewRemoteRoleService(baseURL string) RoleService {
	return &remoteRoleService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteRoleService) ListAll() ([]Role, error) {
	resp, err := s.client.Get(s.baseURL + "/api/roles")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}

	var data struct {
		Roles []Role `json:"roles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Roles, nil
}

func (s *remoteRoleService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]Role, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)

	u := fmt.Sprintf("%s/api/roles?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}

	var data struct {
		Roles []Role `json:"roles"`
		Total int    `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Roles, data.Total, nil
}

func (s *remoteRoleService) FindByID(id string) (*Role, error) {
	u := fmt.Sprintf("%s/api/roles/%s", s.baseURL, url.PathEscape(id))
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}

	var role Role
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *remoteRoleService) CreateRole(id, name string) (*Role, error) {
	payload := map[string]string{
		"id":   id,
		"name": name,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Post(s.baseURL+"/api/roles", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errData map[string]string
		json.NewDecoder(resp.Body).Decode(&errData)
		if msg, ok := errData["error"]; ok {
			return nil, errors.New(msg)
		}
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}

	var role Role
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *remoteRoleService) UpdateRole(id, name string) (*Role, error) {
	payload := map[string]string{
		"name": name,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/roles/%s", s.baseURL, url.PathEscape(id)), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errData map[string]string
		json.NewDecoder(resp.Body).Decode(&errData)
		if msg, ok := errData["error"]; ok {
			return nil, errors.New(msg)
		}
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}

	var role Role
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *remoteRoleService) DeleteRole(id string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/roles/%s", s.baseURL, url.PathEscape(id)), nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errData map[string]string
		json.NewDecoder(resp.Body).Decode(&errData)
		if msg, ok := errData["error"]; ok {
			return errors.New(msg)
		}
		return fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	return nil
}

type remoteAuthService struct {
	client  *http.Client
	baseURL string
}

// NewRemoteAuthService creates a new HTTP-based AuthService client.
func NewRemoteAuthService(baseURL string) AuthService {
	return &remoteAuthService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteAuthService) Login(email, password string) (string, *User, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", nil, err
	}

	resp, err := s.client.Post(s.baseURL+"/api/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errData map[string]string
		json.NewDecoder(resp.Body).Decode(&errData)
		if msg, ok := errData["error"]; ok {
			return "", nil, errors.New(msg)
		}
		return "", nil, fmt.Errorf("remote auth failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Token string `json:"token"`
		User  User   `json:"user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", nil, err
	}
	return result.Token, &result.User, nil
}
