package vehicle

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

type remoteVehicleColorService struct {
	client  *http.Client
	baseURL string
}

// NewRemoteVehicleColorService creates a new HTTP-based VehicleColorService client.
func NewRemoteVehicleColorService(baseURL string) VehicleColorService {
	return &remoteVehicleColorService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteVehicleColorService) ListAll() ([]VehicleColor, error) {
	resp, err := s.client.Get(s.baseURL + "/api/vehicle-colors/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}

	var data struct {
		Colors []VehicleColor `json:"colors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Colors, nil
}

func (s *remoteVehicleColorService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleColor, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)

	u := fmt.Sprintf("%s/api/vehicle-colors?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}

	var data struct {
		Colors []VehicleColor `json:"colors"`
		Total  int            `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Colors, data.Total, nil
}

func (s *remoteVehicleColorService) FindByID(id int64) (*VehicleColor, error) {
	u := fmt.Sprintf("%s/api/vehicle-colors/%d", s.baseURL, id)
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

	var color VehicleColor
	if err := json.NewDecoder(resp.Body).Decode(&color); err != nil {
		return nil, err
	}
	return &color, nil
}

func (s *remoteVehicleColorService) CreateColor(name string, status bool, user string) (*VehicleColor, error) {
	payload := map[string]interface{}{
		"name":   name,
		"status": status,
		"user":   user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Post(s.baseURL+"/api/vehicle-colors", "application/json", bytes.NewBuffer(body))
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

	var color VehicleColor
	if err := json.NewDecoder(resp.Body).Decode(&color); err != nil {
		return nil, err
	}
	return &color, nil
}

func (s *remoteVehicleColorService) UpdateColor(id int64, name string, status bool, user string) (*VehicleColor, error) {
	payload := map[string]interface{}{
		"name":   name,
		"status": status,
		"user":   user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/vehicle-colors/%d", s.baseURL, id), bytes.NewBuffer(body))
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

	var color VehicleColor
	if err := json.NewDecoder(resp.Body).Decode(&color); err != nil {
		return nil, err
	}
	return &color, nil
}

func (s *remoteVehicleColorService) DeleteColor(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/vehicle-colors/%d", s.baseURL, id), nil)
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

type remoteVehicleMakeService struct {
	client  *http.Client
	baseURL string
}

// NewRemoteVehicleMakeService creates a new HTTP-based VehicleMakeService client.
func NewRemoteVehicleMakeService(baseURL string) VehicleMakeService {
	return &remoteVehicleMakeService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteVehicleMakeService) ListAll() ([]VehicleMake, error) {
	resp, err := s.client.Get(s.baseURL + "/api/vehicle-makes/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}

	var data struct {
		Makes []VehicleMake `json:"makes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Makes, nil
}

func (s *remoteVehicleMakeService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleMake, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)

	u := fmt.Sprintf("%s/api/vehicle-makes?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}

	var data struct {
		Makes []VehicleMake `json:"makes"`
		Total int           `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Makes, data.Total, nil
}

func (s *remoteVehicleMakeService) FindByID(id int64) (*VehicleMake, error) {
	u := fmt.Sprintf("%s/api/vehicle-makes/%d", s.baseURL, id)
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

	var make VehicleMake
	if err := json.NewDecoder(resp.Body).Decode(&make); err != nil {
		return nil, err
	}
	return &make, nil
}

func (s *remoteVehicleMakeService) CreateMake(name string, status bool, user string) (*VehicleMake, error) {
	payload := map[string]interface{}{
		"name":   name,
		"status": status,
		"user":   user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Post(s.baseURL+"/api/vehicle-makes", "application/json", bytes.NewBuffer(body))
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

	var make VehicleMake
	if err := json.NewDecoder(resp.Body).Decode(&make); err != nil {
		return nil, err
	}
	return &make, nil
}

func (s *remoteVehicleMakeService) UpdateMake(id int64, name string, status bool, user string) (*VehicleMake, error) {
	payload := map[string]interface{}{
		"name":   name,
		"status": status,
		"user":   user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/vehicle-makes/%d", s.baseURL, id), bytes.NewBuffer(body))
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

	var make VehicleMake
	if err := json.NewDecoder(resp.Body).Decode(&make); err != nil {
		return nil, err
	}
	return &make, nil
}

func (s *remoteVehicleMakeService) DeleteMake(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/vehicle-makes/%d", s.baseURL, id), nil)
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
