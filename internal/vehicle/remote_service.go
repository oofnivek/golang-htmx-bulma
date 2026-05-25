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
