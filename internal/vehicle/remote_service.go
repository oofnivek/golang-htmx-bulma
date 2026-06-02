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

	"golang-htmx-bulma/internal/pkg/status"
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

// remoteFuelTypeService provides HTTP-based client for FuelTypeService.

type remoteFuelTypeService struct {
	client  *http.Client
	baseURL string
}

func NewRemoteFuelTypeService(baseURL string) FuelTypeService {
	return &remoteFuelTypeService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteFuelTypeService) ListAll() ([]FuelType, error) {
	resp, err := s.client.Get(s.baseURL + "/api/fuel-types/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Fuels []FuelType `json:"fuels"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Fuels, nil
}

func (s *remoteFuelTypeService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]FuelType, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)
	u := fmt.Sprintf("%s/api/fuel-types?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Fuels []FuelType `json:"fuels"`
		Total int        `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Fuels, data.Total, nil
}

func (s *remoteFuelTypeService) FindByID(id int64) (*FuelType, error) {
	u := fmt.Sprintf("%s/api/fuel-types/%d", s.baseURL, id)
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
	var f FuelType
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *remoteFuelTypeService) CreateFuelType(name string, status bool, user string) (*FuelType, error) {
	payload := map[string]interface{}{"name": name, "status": status, "user": user}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(s.baseURL+"/api/fuel-types", "application/json", bytes.NewBuffer(body))
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
	var f FuelType
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *remoteFuelTypeService) UpdateFuelType(id int64, name string, status bool, user string) (*FuelType, error) {
	payload := map[string]interface{}{"name": name, "status": status, "user": user}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/fuel-types/%d", s.baseURL, id), bytes.NewBuffer(body))
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
	var f FuelType
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *remoteFuelTypeService) DeleteFuelType(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/fuel-types/%d", s.baseURL, id), nil)
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

// remoteVehicleStatusService provides HTTP-based client for VehicleStatusService.

type remoteVehicleStatusService struct {
	client  *http.Client
	baseURL string
}

func NewRemoteVehicleStatusService(baseURL string) VehicleStatusService {
	return &remoteVehicleStatusService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteVehicleStatusService) ListAll() ([]VehicleStatus, error) {
	resp, err := s.client.Get(s.baseURL + "/api/vehicle-statuses/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Statuses []VehicleStatus `json:"statuses"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Statuses, nil
}

func (s *remoteVehicleStatusService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleStatus, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)
	u := fmt.Sprintf("%s/api/vehicle-statuses?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Statuses []VehicleStatus `json:"statuses"`
		Total    int             `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Statuses, data.Total, nil
}

func (s *remoteVehicleStatusService) FindByID(id int64) (*VehicleStatus, error) {
	u := fmt.Sprintf("%s/api/vehicle-statuses/%d", s.baseURL, id)
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
	var vs VehicleStatus
	if err := json.NewDecoder(resp.Body).Decode(&vs); err != nil {
		return nil, err
	}
	return &vs, nil
}

func (s *remoteVehicleStatusService) CreateStatus(substatus string, isActive bool) (*VehicleStatus, error) {
	payload := map[string]interface{}{"substatus": substatus, "is_active": isActive}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(s.baseURL+"/api/vehicle-statuses", "application/json", bytes.NewBuffer(body))
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
	var vs VehicleStatus
	if err := json.NewDecoder(resp.Body).Decode(&vs); err != nil {
		return nil, err
	}
	return &vs, nil
}

func (s *remoteVehicleStatusService) UpdateStatus(id int64, substatus string, isActive bool) (*VehicleStatus, error) {
	payload := map[string]interface{}{"substatus": substatus, "is_active": isActive}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/vehicle-statuses/%d", s.baseURL, id), bytes.NewBuffer(body))
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
	var vs VehicleStatus
	if err := json.NewDecoder(resp.Body).Decode(&vs); err != nil {
		return nil, err
	}
	return &vs, nil
}

func (s *remoteVehicleStatusService) DeleteStatus(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/vehicle-statuses/%d", s.baseURL, id), nil)
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

// RemoteVehicleTypeService provides HTTP-based client for VehicleTypeService.

type remoteVehicleTypeService struct {
    client  *http.Client
    baseURL string
}

// NewRemoteVehicleTypeService creates a new HTTP-based VehicleTypeService client.
func NewRemoteVehicleTypeService(baseURL string) VehicleTypeService {
    return &remoteVehicleTypeService{
        client:  &http.Client{Timeout: 10 * time.Second},
        baseURL: strings.TrimRight(baseURL, "/"),
    }
}

func (s *remoteVehicleTypeService) ListAll() ([]VehicleType, error) {
    resp, err := s.client.Get(s.baseURL + "/api/vehicle-types/all")
    if err != nil {
        return nil, fmt.Errorf("remote service call failed: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
    }
    var data struct { Types []VehicleType `json:"types"` }
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, err
    }
    return data.Types, nil
}

func (s *remoteVehicleTypeService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleType, int, error) {
    q := url.Values{}
    q.Set("page", strconv.Itoa(page))
    q.Set("pageSize", strconv.Itoa(pageSize))
    q.Set("sortBy", sortBy)
    q.Set("sortOrder", sortOrder)
    u := fmt.Sprintf("%s/api/vehicle-types?%s", s.baseURL, q.Encode())
    resp, err := s.client.Get(u)
    if err != nil {
        return nil, 0, fmt.Errorf("remote service call failed: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
    }
    var data struct {
        Types []VehicleType `json:"types"`
        Total int           `json:"total"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, 0, err
    }
    return data.Types, data.Total, nil
}

func (s *remoteVehicleTypeService) FindByID(id int64) (*VehicleType, error) {
    u := fmt.Sprintf("%s/api/vehicle-types/%d", s.baseURL, id)
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
    var vt VehicleType
    if err := json.NewDecoder(resp.Body).Decode(&vt); err != nil {
        return nil, err
    }
    return &vt, nil
}

func (s *remoteVehicleTypeService) Create(name string, status bool, user string) (*VehicleType, error) {
    payload := map[string]interface{}{ "name": name, "status": status, "user": user }
    body, err := json.Marshal(payload)
    if err != nil { return nil, err }
    resp, err := s.client.Post(s.baseURL+"/api/vehicle-types", "application/json", bytes.NewBuffer(body))
    if err != nil { return nil, fmt.Errorf("remote service call failed: %w", err) }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
        var errData map[string]string
        json.NewDecoder(resp.Body).Decode(&errData)
        if msg, ok := errData["error"]; ok { return nil, errors.New(msg) }
        return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
    }
    var vt VehicleType
    if err := json.NewDecoder(resp.Body).Decode(&vt); err != nil { return nil, err }
    return &vt, nil
}

func (s *remoteVehicleTypeService) Update(id int64, name string, status bool, user string) (*VehicleType, error) {
    payload := map[string]interface{}{ "name": name, "status": status, "user": user }
    body, err := json.Marshal(payload)
    if err != nil { return nil, err }
    req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/vehicle-types/%d", s.baseURL, id), bytes.NewBuffer(body))
    if err != nil { return nil, err }
    req.Header.Set("Content-Type", "application/json")
    resp, err := s.client.Do(req)
    if err != nil { return nil, fmt.Errorf("remote service call failed: %w", err) }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        var errData map[string]string
        json.NewDecoder(resp.Body).Decode(&errData)
        if msg, ok := errData["error"]; ok { return nil, errors.New(msg) }
        return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
    }
    var vt VehicleType
    if err := json.NewDecoder(resp.Body).Decode(&vt); err != nil { return nil, err }
    return &vt, nil
}

func (s *remoteVehicleTypeService) Delete(id int64) error {
    req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/vehicle-types/%d", s.baseURL, id), nil)
    if err != nil { return err }
    resp, err := s.client.Do(req)
    if err != nil { return fmt.Errorf("remote service call failed: %w", err) }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        var errData map[string]string
        json.NewDecoder(resp.Body).Decode(&errData)
        if msg, ok := errData["error"]; ok { return errors.New(msg) }
        return fmt.Errorf("remote service returned status: %d", resp.StatusCode)
    }
    return nil
}

// remoteVehicleFuelService provides HTTP-based client for VehicleFuelService.

type remoteVehicleFuelService struct {
	client  *http.Client
	baseURL string
}

func NewRemoteVehicleFuelService(baseURL string) VehicleFuelService {
	return &remoteVehicleFuelService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteVehicleFuelService) ListAll() ([]VehicleFuel, error) {
	resp, err := s.client.Get(s.baseURL + "/api/vehicle-fuels/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Fuels []VehicleFuel `json:"fuels"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Fuels, nil
}

func (s *remoteVehicleFuelService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleFuel, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)
	u := fmt.Sprintf("%s/api/vehicle-fuels?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Fuels []VehicleFuel `json:"fuels"`
		Total int           `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Fuels, data.Total, nil
}

func (s *remoteVehicleFuelService) FindByID(id int64) (*VehicleFuel, error) {
	u := fmt.Sprintf("%s/api/vehicle-fuels/%d", s.baseURL, id)
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
	var f VehicleFuel
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *remoteVehicleFuelService) CreateFuel(vehicleMakeID, vehicleModelID, fuelTypeID int64, fuelTankSize, fuelConsumption float64, status bool, user string) (*VehicleFuel, error) {
	payload := map[string]interface{}{
		"vehicle_make_id":  vehicleMakeID,
		"vehicle_model_id": vehicleModelID,
		"fuel_type_id":     fuelTypeID,
		"fuel_tank_size":   fuelTankSize,
		"fuel_consumption": fuelConsumption,
		"status":           status,
		"user":             user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(s.baseURL+"/api/vehicle-fuels", "application/json", bytes.NewBuffer(body))
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
	var f VehicleFuel
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *remoteVehicleFuelService) UpdateFuel(id, vehicleMakeID, vehicleModelID, fuelTypeID int64, fuelTankSize, fuelConsumption float64, status bool, user string) (*VehicleFuel, error) {
	payload := map[string]interface{}{
		"vehicle_make_id":  vehicleMakeID,
		"vehicle_model_id": vehicleModelID,
		"fuel_type_id":     fuelTypeID,
		"fuel_tank_size":   fuelTankSize,
		"fuel_consumption": fuelConsumption,
		"status":           status,
		"user":             user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/vehicle-fuels/%d", s.baseURL, id), bytes.NewBuffer(body))
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
	var f VehicleFuel
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *remoteVehicleFuelService) DeleteFuel(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/vehicle-fuels/%d", s.baseURL, id), nil)
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

// remoteFuelCompanyService provides HTTP-based client for FuelCompanyService.

type remoteFuelCompanyService struct {
	client  *http.Client
	baseURL string
}

func NewRemoteFuelCompanyService(baseURL string) FuelCompanyService {
	return &remoteFuelCompanyService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteFuelCompanyService) ListAll() ([]FuelCompany, error) {
	resp, err := s.client.Get(s.baseURL + "/api/fuel-companies/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Companies []FuelCompany `json:"companies"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Companies, nil
}

func (s *remoteFuelCompanyService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]FuelCompany, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)
	u := fmt.Sprintf("%s/api/fuel-companies?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Companies []FuelCompany `json:"companies"`
		Total     int           `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Companies, data.Total, nil
}

func (s *remoteFuelCompanyService) FindByID(id int64) (*FuelCompany, error) {
	u := fmt.Sprintf("%s/api/fuel-companies/%d", s.baseURL, id)
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
	var c FuelCompany
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *remoteFuelCompanyService) CreateFuelCompany(name string, status bool, user string) (*FuelCompany, error) {
	payload := map[string]interface{}{"name": name, "status": status, "user": user}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(s.baseURL+"/api/fuel-companies", "application/json", bytes.NewBuffer(body))
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
	var c FuelCompany
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *remoteFuelCompanyService) UpdateFuelCompany(id int64, name string, status bool, user string) (*FuelCompany, error) {
	payload := map[string]interface{}{"name": name, "status": status, "user": user}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/fuel-companies/%d", s.baseURL, id), bytes.NewBuffer(body))
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
	var c FuelCompany
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *remoteFuelCompanyService) DeleteFuelCompany(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/fuel-companies/%d", s.baseURL, id), nil)
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

// remoteVehicleModelService provides HTTP-based client for VehicleModelService.

type remoteVehicleModelService struct {
	client  *http.Client
	baseURL string
}

func NewRemoteVehicleModelService(baseURL string) VehicleModelService {
	return &remoteVehicleModelService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteVehicleModelService) ListAll() ([]VehicleModel, error) {
	resp, err := s.client.Get(s.baseURL + "/api/vehicle-models/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Models []VehicleModel `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Models, nil
}

func (s *remoteVehicleModelService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]VehicleModel, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)
	u := fmt.Sprintf("%s/api/vehicle-models?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Models []VehicleModel `json:"models"`
		Total  int            `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Models, data.Total, nil
}

func (s *remoteVehicleModelService) FindByID(id int64) (*VehicleModel, error) {
	u := fmt.Sprintf("%s/api/vehicle-models/%d", s.baseURL, id)
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
	var m VehicleModel
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *remoteVehicleModelService) CreateModel(vehicleTypeID, vehicleMakeID int64, name string, status bool, user string) (*VehicleModel, error) {
	payload := map[string]interface{}{
		"vehicle_type_id": vehicleTypeID,
		"vehicle_make_id": vehicleMakeID,
		"name":            name,
		"status":          status,
		"user":            user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(s.baseURL+"/api/vehicle-models", "application/json", bytes.NewBuffer(body))
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
	var m VehicleModel
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *remoteVehicleModelService) UpdateModel(id, vehicleTypeID, vehicleMakeID int64, name string, status bool, user string) (*VehicleModel, error) {
	payload := map[string]interface{}{
		"vehicle_type_id": vehicleTypeID,
		"vehicle_make_id": vehicleMakeID,
		"name":            name,
		"status":          status,
		"user":            user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/vehicle-models/%d", s.baseURL, id), bytes.NewBuffer(body))
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
	var m VehicleModel
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *remoteVehicleModelService) DeleteModel(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/vehicle-models/%d", s.baseURL, id), nil)
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

// remoteCarAssetOwnerService provides HTTP-based client for CarAssetOwnerService.

type remoteCarAssetOwnerService struct {
	client  *http.Client
	baseURL string
}

func NewRemoteCarAssetOwnerService(baseURL string) CarAssetOwnerService {
	return &remoteCarAssetOwnerService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteCarAssetOwnerService) ListAll() ([]CarAssetOwner, error) {
	resp, err := s.client.Get(s.baseURL + "/api/car-asset-owners/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Owners []CarAssetOwner `json:"owners"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Owners, nil
}

func (s *remoteCarAssetOwnerService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CarAssetOwner, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)
	u := fmt.Sprintf("%s/api/car-asset-owners?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Owners []CarAssetOwner `json:"owners"`
		Total  int             `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Owners, data.Total, nil
}

func (s *remoteCarAssetOwnerService) FindByID(id int64) (*CarAssetOwner, error) {
	u := fmt.Sprintf("%s/api/car-asset-owners/%d", s.baseURL, id)
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
	var o CarAssetOwner
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *remoteCarAssetOwnerService) CreateOwner(name string, status bool, user string) (*CarAssetOwner, error) {
	payload := map[string]interface{}{"name": name, "status": status, "user": user}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(s.baseURL+"/api/car-asset-owners", "application/json", bytes.NewBuffer(body))
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
	var o CarAssetOwner
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *remoteCarAssetOwnerService) UpdateOwner(id int64, name string, status bool, user string) (*CarAssetOwner, error) {
	payload := map[string]interface{}{"name": name, "status": status, "user": user}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/car-asset-owners/%d", s.baseURL, id), bytes.NewBuffer(body))
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
	var o CarAssetOwner
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *remoteCarAssetOwnerService) DeleteOwner(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/car-asset-owners/%d", s.baseURL, id), nil)
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

// remoteCarParkOwnerService provides HTTP-based client for CarParkOwnerService.

type remoteCarParkOwnerService struct {
	client  *http.Client
	baseURL string
}

func NewRemoteCarParkOwnerService(baseURL string) CarParkOwnerService {
	return &remoteCarParkOwnerService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteCarParkOwnerService) ListAll() ([]CarParkOwner, error) {
	resp, err := s.client.Get(s.baseURL + "/api/car-park-owners/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Owners []CarParkOwner `json:"owners"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Owners, nil
}

func (s *remoteCarParkOwnerService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CarParkOwner, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)
	u := fmt.Sprintf("%s/api/car-park-owners?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Owners []CarParkOwner `json:"owners"`
		Total  int            `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Owners, data.Total, nil
}

func (s *remoteCarParkOwnerService) FindByID(id int64) (*CarParkOwner, error) {
	u := fmt.Sprintf("%s/api/car-park-owners/%d", s.baseURL, id)
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
	var o CarParkOwner
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *remoteCarParkOwnerService) CreateOwner(name string, status bool, user string) (*CarParkOwner, error) {
	payload := map[string]interface{}{"name": name, "status": status, "user": user}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(s.baseURL+"/api/car-park-owners", "application/json", bytes.NewBuffer(body))
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
	var o CarParkOwner
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *remoteCarParkOwnerService) UpdateOwner(id int64, name string, status bool, user string) (*CarParkOwner, error) {
	payload := map[string]interface{}{"name": name, "status": status, "user": user}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/car-park-owners/%d", s.baseURL, id), bytes.NewBuffer(body))
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
	var o CarParkOwner
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *remoteCarParkOwnerService) DeleteOwner(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/car-park-owners/%d", s.baseURL, id), nil)
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

// remoteCarParkService provides HTTP-based client for CarParkService.

type remoteCarParkService struct {
	client  *http.Client
	baseURL string
}

func NewRemoteCarParkService(baseURL string) CarParkService {
	return &remoteCarParkService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteCarParkService) ListAll() ([]CarPark, error) {
	resp, err := s.client.Get(s.baseURL + "/api/car-parks/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Parks []CarPark `json:"parks"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Parks, nil
}

func (s *remoteCarParkService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CarPark, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)
	u := fmt.Sprintf("%s/api/car-parks?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Parks []CarPark `json:"parks"`
		Total int       `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Parks, data.Total, nil
}

func (s *remoteCarParkService) FindByID(id int64) (*CarPark, error) {
	u := fmt.Sprintf("%s/api/car-parks/%d", s.baseURL, id)
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
	var cp CarPark
	if err := json.NewDecoder(resp.Body).Decode(&cp); err != nil {
		return nil, err
	}
	return &cp, nil
}

func (s *remoteCarParkService) CreateCarPark(name string, description *string, postalCode, address string, latitude, longitude float64, carParkOwnerID int64, activeFrom, activeTo *time.Time, status bool, user string) (*CarPark, error) {
	payload := map[string]interface{}{
		"name":              name,
		"description":       description,
		"postal_code":       postalCode,
		"address":           address,
		"latitude":          latitude,
		"longitude":         longitude,
		"car_park_owner_id": carParkOwnerID,
		"active_from":       activeFrom,
		"active_to":         activeTo,
		"status":            status,
		"user":              user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(s.baseURL+"/api/car-parks", "application/json", bytes.NewBuffer(body))
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
	var cp CarPark
	if err := json.NewDecoder(resp.Body).Decode(&cp); err != nil {
		return nil, err
	}
	return &cp, nil
}

func (s *remoteCarParkService) UpdateCarPark(id int64, name string, description *string, postalCode, address string, latitude, longitude float64, carParkOwnerID int64, activeFrom, activeTo *time.Time, status bool, user string) (*CarPark, error) {
	payload := map[string]interface{}{
		"name":              name,
		"description":       description,
		"postal_code":       postalCode,
		"address":           address,
		"latitude":          latitude,
		"longitude":         longitude,
		"car_park_owner_id": carParkOwnerID,
		"active_from":       activeFrom,
		"active_to":         activeTo,
		"status":            status,
		"user":              user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/car-parks/%d", s.baseURL, id), bytes.NewBuffer(body))
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
	var cp CarPark
	if err := json.NewDecoder(resp.Body).Decode(&cp); err != nil {
		return nil, err
	}
	return &cp, nil
}

func (s *remoteCarParkService) DeleteCarPark(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/car-parks/%d", s.baseURL, id), nil)
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

// remoteCarParkLotService provides HTTP-based client for CarParkLotService.

type remoteCarParkLotService struct {
	client  *http.Client
	baseURL string
}

func NewRemoteCarParkLotService(baseURL string) CarParkLotService {
	return &remoteCarParkLotService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteCarParkLotService) ListAll() ([]CarParkLot, error) {
	resp, err := s.client.Get(s.baseURL + "/api/car-park-lots/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Lots []CarParkLot `json:"lots"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Lots, nil
}

func (s *remoteCarParkLotService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]CarParkLot, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)
	u := fmt.Sprintf("%s/api/car-park-lots?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Lots  []CarParkLot `json:"lots"`
		Total int          `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Lots, data.Total, nil
}

func (s *remoteCarParkLotService) FindByID(id int64) (*CarParkLot, error) {
	u := fmt.Sprintf("%s/api/car-park-lots/%d", s.baseURL, id)
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
	var l CarParkLot
	if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
		return nil, err
	}
	return &l, nil
}

func (s *remoteCarParkLotService) CreateCarParkLot(carParkID int64, lotNumber, level string, status bool, user string) (*CarParkLot, error) {
	payload := map[string]interface{}{
		"car_park_id": carParkID,
		"lot_number":  lotNumber,
		"level":       level,
		"status":      status,
		"user":        user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(s.baseURL+"/api/car-park-lots", "application/json", bytes.NewBuffer(body))
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
	var l CarParkLot
	if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
		return nil, err
	}
	return &l, nil
}

func (s *remoteCarParkLotService) UpdateCarParkLot(id, carParkID int64, lotNumber, level string, status bool, user string) (*CarParkLot, error) {
	payload := map[string]interface{}{
		"car_park_id": carParkID,
		"lot_number":  lotNumber,
		"level":       level,
		"status":      status,
		"user":        user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/car-park-lots/%d", s.baseURL, id), bytes.NewBuffer(body))
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
	var l CarParkLot
	if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
		return nil, err
	}
	return &l, nil
}

func (s *remoteCarParkLotService) DeleteCarParkLot(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/car-park-lots/%d", s.baseURL, id), nil)
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

// remoteEstateService provides HTTP-based client for EstateService.

type remoteEstateService struct {
	client  *http.Client
	baseURL string
}

func NewRemoteEstateService(baseURL string) EstateService {
	return &remoteEstateService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteEstateService) ListAll() ([]Estate, error) {
	resp, err := s.client.Get(s.baseURL + "/api/estates/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Estates []Estate `json:"estates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Estates, nil
}

func (s *remoteEstateService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]Estate, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)
	u := fmt.Sprintf("%s/api/estates?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Estates []Estate `json:"estates"`
		Total   int      `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Estates, data.Total, nil
}

func (s *remoteEstateService) FindByID(id int64) (*Estate, error) {
	u := fmt.Sprintf("%s/api/estates/%d", s.baseURL, id)
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
	var e Estate
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *remoteEstateService) CreateEstate(name string) (*Estate, error) {
	payload := map[string]interface{}{"name": name}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(s.baseURL+"/api/estates", "application/json", bytes.NewBuffer(body))
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
	var e Estate
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *remoteEstateService) UpdateEstate(id int64, name string) (*Estate, error) {
	payload := map[string]interface{}{"name": name}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/estates/%d", s.baseURL, id), bytes.NewBuffer(body))
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
	var e Estate
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *remoteEstateService) DeleteEstate(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/estates/%d", s.baseURL, id), nil)
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

// remoteRegionalInfoService provides HTTP-based client for RegionalInfoService.
type remoteRegionalInfoService struct {
	baseURL string
	client  *http.Client
}

func NewRemoteRegionalInfoService(baseURL string) RegionalInfoService {
	return &remoteRegionalInfoService{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (s *remoteRegionalInfoService) ListAll() ([]RegionalInfo, error) {
	resp, err := s.client.Get(s.baseURL + "/api/regional-infos/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	var data struct {
		RegionalInfos []RegionalInfo `json:"regional_infos"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.RegionalInfos, nil
}

func (s *remoteRegionalInfoService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]RegionalInfo, int, error) {
	q := url.Values{}
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("pageSize", fmt.Sprintf("%d", pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)
	u := fmt.Sprintf("%s/api/regional-infos?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	var data struct {
		RegionalInfos []RegionalInfo `json:"regional_infos"`
		Total         int            `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.RegionalInfos, data.Total, nil
}

func (s *remoteRegionalInfoService) FindByID(postalCode string) (*RegionalInfo, error) {
	u := fmt.Sprintf("%s/api/regional-infos/%s", s.baseURL, postalCode)
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	var ri RegionalInfo
	if err := json.NewDecoder(resp.Body).Decode(&ri); err != nil {
		return nil, err
	}
	return &ri, nil
}

func (s *remoteRegionalInfoService) CreateRegionalInfo(postalCode, region string, estateID int64) (*RegionalInfo, error) {
	body, err := json.Marshal(map[string]any{
		"postal_code": postalCode,
		"region":      region,
		"estate_id":   estateID,
	})
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(s.baseURL+"/api/regional-infos", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		var errData map[string]string
		json.NewDecoder(resp.Body).Decode(&errData)
		if msg, ok := errData["error"]; ok {
			return nil, errors.New(msg)
		}
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var ri RegionalInfo
	if err := json.NewDecoder(resp.Body).Decode(&ri); err != nil {
		return nil, err
	}
	return &ri, nil
}

func (s *remoteRegionalInfoService) UpdateRegionalInfo(postalCode, region string, estateID int64) (*RegionalInfo, error) {
	body, err := json.Marshal(map[string]any{
		"region":    region,
		"estate_id": estateID,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/regional-infos/%s", s.baseURL, postalCode), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		var errData map[string]string
		json.NewDecoder(resp.Body).Decode(&errData)
		if msg, ok := errData["error"]; ok {
			return nil, errors.New(msg)
		}
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var ri RegionalInfo
	if err := json.NewDecoder(resp.Body).Decode(&ri); err != nil {
		return nil, err
	}
	return &ri, nil
}

func (s *remoteRegionalInfoService) DeleteRegionalInfo(postalCode string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/regional-infos/%s", s.baseURL, postalCode), nil)
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

// remoteVehicleService provides HTTP-based client for VehicleService.

type remoteVehicleService struct {
	client  *http.Client
	baseURL string
}

func NewRemoteVehicleService(baseURL string) VehicleService {
	return &remoteVehicleService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteVehicleService) ListAll() ([]Vehicle, error) {
	resp, err := s.client.Get(s.baseURL + "/api/vehicles/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Vehicles []Vehicle `json:"vehicles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Vehicles, nil
}

func (s *remoteVehicleService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]Vehicle, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)
	u := fmt.Sprintf("%s/api/vehicles?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Vehicles []Vehicle `json:"vehicles"`
		Total    int       `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Vehicles, data.Total, nil
}

func (s *remoteVehicleService) FindByID(id int64) (*Vehicle, error) {
	u := fmt.Sprintf("%s/api/vehicles/%d", s.baseURL, id)
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
	var v Vehicle
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *remoteVehicleService) CreateVehicle(
	vehicleMakeID, vehicleModelID, vehicleTypeID, fuelTypeID, vehicleColorID int64,
	description, plateNumber, iuNumber, chassisNumber, engineNumber, bootSpace *string,
	numSeats int,
	carParkID, assetOwnerID, vehicleStatusID int64,
	lastServiceDate, lastCleanedDate, activeFrom, activeTo *time.Time,
	lastServiceMileage, currentMileage, currentFuelLevel *int,
	user string,
) (*Vehicle, error) {
	payload := map[string]interface{}{
		"vehicle_make_id":      vehicleMakeID,
		"vehicle_model_id":     vehicleModelID,
		"vehicle_type_id":      vehicleTypeID,
		"fuel_type_id":         fuelTypeID,
		"vehicle_color_id":     vehicleColorID,
		"description":          description,
		"plate_number":         plateNumber,
		"iu_number":            iuNumber,
		"chassis_number":       chassisNumber,
		"engine_number":        engineNumber,
		"num_seats":            numSeats,
		"boot_space":           bootSpace,
		"car_park_id":          carParkID,
		"asset_owner_id":       assetOwnerID,
		"vehicle_status_id":    vehicleStatusID,
		"last_service_date":    lastServiceDate,
		"last_cleaned_date":    lastCleanedDate,
		"last_service_mileage": lastServiceMileage,
		"current_mileage":      currentMileage,
		"current_fuel_level":   currentFuelLevel,
		"active_from":          activeFrom,
		"active_to":            activeTo,
		"user":                 user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(s.baseURL+"/api/vehicles", "application/json", bytes.NewBuffer(body))
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
	var v Vehicle
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *remoteVehicleService) UpdateVehicle(
	id int64,
	vehicleMakeID, vehicleModelID, vehicleTypeID, fuelTypeID, vehicleColorID int64,
	description, plateNumber, iuNumber, chassisNumber, engineNumber, bootSpace *string,
	numSeats int,
	carParkID, assetOwnerID, vehicleStatusID int64,
	lastServiceDate, lastCleanedDate, activeFrom, activeTo *time.Time,
	lastServiceMileage, currentMileage, currentFuelLevel *int,
	user string,
) (*Vehicle, error) {
	payload := map[string]interface{}{
		"vehicle_make_id":      vehicleMakeID,
		"vehicle_model_id":     vehicleModelID,
		"vehicle_type_id":      vehicleTypeID,
		"fuel_type_id":         fuelTypeID,
		"vehicle_color_id":     vehicleColorID,
		"description":          description,
		"plate_number":         plateNumber,
		"iu_number":            iuNumber,
		"chassis_number":       chassisNumber,
		"engine_number":        engineNumber,
		"num_seats":            numSeats,
		"boot_space":           bootSpace,
		"car_park_id":          carParkID,
		"asset_owner_id":       assetOwnerID,
		"vehicle_status_id":    vehicleStatusID,
		"last_service_date":    lastServiceDate,
		"last_cleaned_date":    lastCleanedDate,
		"last_service_mileage": lastServiceMileage,
		"current_mileage":      currentMileage,
		"current_fuel_level":   currentFuelLevel,
		"active_from":          activeFrom,
		"active_to":            activeTo,
		"user":                 user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/vehicles/%d", s.baseURL, id), bytes.NewBuffer(body))
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
	var v Vehicle
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *remoteVehicleService) DeleteVehicle(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/vehicles/%d", s.baseURL, id), nil)
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

// remoteFuelCardService provides HTTP-based client for FuelCardService.

type remoteFuelCardService struct {
	client  *http.Client
	baseURL string
}

func NewRemoteFuelCardService(baseURL string) FuelCardService {
	return &remoteFuelCardService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteFuelCardService) ListAll() ([]FuelCard, error) {
	resp, err := s.client.Get(s.baseURL + "/api/fuel-cards/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Cards []FuelCard `json:"cards"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Cards, nil
}

func (s *remoteFuelCardService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]FuelCard, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)
	u := fmt.Sprintf("%s/api/fuel-cards?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Cards []FuelCard `json:"cards"`
		Total int        `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Cards, data.Total, nil
}

func (s *remoteFuelCardService) FindByID(id int64) (*FuelCard, error) {
	u := fmt.Sprintf("%s/api/fuel-cards/%d", s.baseURL, id)
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
	var fc FuelCard
	if err := json.NewDecoder(resp.Body).Decode(&fc); err != nil {
		return nil, err
	}
	return &fc, nil
}

func (s *remoteFuelCardService) CreateFuelCard(cardNo string, fuelCompanyID int64, pinNumber string, vehicleID *int64, status bool, user string) (*FuelCard, error) {
	payload := map[string]interface{}{
		"card_no":         cardNo,
		"fuel_company_id": fuelCompanyID,
		"pin_number":      pinNumber,
		"vehicle_id":      vehicleID,
		"status":          status,
		"user":            user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(s.baseURL+"/api/fuel-cards", "application/json", bytes.NewBuffer(body))
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
	var fc FuelCard
	if err := json.NewDecoder(resp.Body).Decode(&fc); err != nil {
		return nil, err
	}
	return &fc, nil
}

func (s *remoteFuelCardService) UpdateFuelCard(id int64, cardNo string, fuelCompanyID int64, pinNumber string, vehicleID *int64, status bool, user string) (*FuelCard, error) {
	payload := map[string]interface{}{
		"card_no":         cardNo,
		"fuel_company_id": fuelCompanyID,
		"pin_number":      pinNumber,
		"vehicle_id":      vehicleID,
		"status":          status,
		"user":            user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/fuel-cards/%d", s.baseURL, id), bytes.NewBuffer(body))
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
	var fc FuelCard
	if err := json.NewDecoder(resp.Body).Decode(&fc); err != nil {
		return nil, err
	}
	return &fc, nil
}

func (s *remoteFuelCardService) DeleteFuelCard(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/fuel-cards/%d", s.baseURL, id), nil)
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

// remoteCondoService provides HTTP-based client for CondoService.

type remoteCondoService struct {
	client  *http.Client
	baseURL string
}

func NewRemoteCondoService(baseURL string) CondoService {
	return &remoteCondoService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *remoteCondoService) ListAll() ([]Condo, error) {
	resp, err := s.client.Get(s.baseURL + "/api/condos/all")
	if err != nil {
		return nil, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Condos []Condo `json:"condos"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Condos, nil
}

func (s *remoteCondoService) ListPaged(page, pageSize int, sortBy, sortOrder string) ([]Condo, int, error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("sortBy", sortBy)
	q.Set("sortOrder", sortOrder)
	u := fmt.Sprintf("%s/api/condos?%s", s.baseURL, q.Encode())
	resp, err := s.client.Get(u)
	if err != nil {
		return nil, 0, fmt.Errorf("remote service call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("remote service returned status: %d", resp.StatusCode)
	}
	var data struct {
		Condos []Condo `json:"condos"`
		Total  int     `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, 0, err
	}
	return data.Condos, data.Total, nil
}

func (s *remoteCondoService) FindByID(id int64) (*Condo, error) {
	u := fmt.Sprintf("%s/api/condos/%d", s.baseURL, id)
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
	var c Condo
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *remoteCondoService) CreateCondo(name string, st status.Status, mcstNumber, mcstEmail, address, user string) (*Condo, error) {
	payload := map[string]interface{}{
		"name":        name,
		"status":      st,
		"mcst_number": mcstNumber,
		"mcst_email":  mcstEmail,
		"address":     address,
		"user":        user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(s.baseURL+"/api/condos", "application/json", bytes.NewBuffer(body))
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
	var c Condo
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *remoteCondoService) UpdateCondo(id int64, name string, st status.Status, mcstNumber, mcstEmail, address, user string) (*Condo, error) {
	payload := map[string]interface{}{
		"name":        name,
		"status":      st,
		"mcst_number": mcstNumber,
		"mcst_email":  mcstEmail,
		"address":     address,
		"user":        user,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/condos/%d", s.baseURL, id), bytes.NewBuffer(body))
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
	var c Condo
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *remoteCondoService) DeleteCondo(id int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/condos/%d", s.baseURL, id), nil)
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
