package vehicle

import (
	"errors"
	"testing"
)

func TestCreateColor(t *testing.T) {
	mockRepo := &MockVehicleColorRepository{
		CreateFn: func(color *VehicleColor) error {
			color.ID = 1
			return nil
		},
	}

	svc := NewVehicleColorService(mockRepo)

	color, err := svc.CreateColor("Red", true, "test-user")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if color.ID != 1 {
		t.Errorf("Expected ID 1, got %d", color.ID)
	}

	if color.Name != "Red" {
		t.Errorf("Expected name 'Red', got %s", color.Name)
	}
}

func TestFindByID(t *testing.T) {
	mockRepo := &MockVehicleColorRepository{
		GetByIDFn: func(id int64) (*VehicleColor, error) {
			if id == 1 {
				return &VehicleColor{ID: 1, Name: "Blue"}, nil
			}
			return nil, nil
		},
	}

	svc := NewVehicleColorService(mockRepo)

	t.Run("Found", func(t *testing.T) {
		color, err := svc.FindByID(1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if color == nil || color.Name != "Blue" {
			t.Errorf("Expected 'Blue', got %v", color)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		color, err := svc.FindByID(99)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if color != nil {
			t.Errorf("Expected nil, got %v", color)
		}
	})
}

func TestListAll(t *testing.T) {
	mockRepo := &MockVehicleColorRepository{
		GetAllFn: func() ([]VehicleColor, error) {
			return []VehicleColor{
				{ID: 1, Name: "Red"},
				{ID: 2, Name: "Blue"},
			}, nil
		},
	}

	svc := NewVehicleColorService(mockRepo)
	colors, err := svc.ListAll()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(colors) != 2 {
		t.Errorf("Expected 2 colors, got %d", len(colors))
	}
}

func TestUpdateColor(t *testing.T) {
	mockRepo := &MockVehicleColorRepository{
		GetByIDFn: func(id int64) (*VehicleColor, error) {
			if id == 1 {
				return &VehicleColor{ID: 1, Name: "Old Name"}, nil
			}
			return nil, nil
		},
		UpdateFn: func(color *VehicleColor) error {
			return nil
		},
	}

	svc := NewVehicleColorService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		color, err := svc.UpdateColor(1, "New Name", false, "admin")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if color.Name != "New Name" {
			t.Errorf("Expected 'New Name', got %s", color.Name)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		color, err := svc.UpdateColor(99, "New Name", false, "admin")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if color != nil {
			t.Errorf("Expected nil for non-existent color")
		}
	})
}

func TestDeleteColor(t *testing.T) {
	mockRepo := &MockVehicleColorRepository{
		DeleteFn: func(id int64) error {
			return nil
		},
	}

	svc := NewVehicleColorService(mockRepo)
	err := svc.DeleteColor(1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestServiceErrors(t *testing.T) {
	errRepo := &MockVehicleColorRepository{
		GetAllFn: func() ([]VehicleColor, error) {
			return nil, errors.New("db error")
		},
		CreateFn: func(color *VehicleColor) error {
			return errors.New("db error")
		},
		UpdateFn: func(color *VehicleColor) error {
			return errors.New("db error")
		},
		GetByIDFn: func(id int64) (*VehicleColor, error) {
			return nil, errors.New("db error")
		},
		DeleteFn: func(id int64) error {
			return errors.New("db error")
		},
	}

	svc := NewVehicleColorService(errRepo)

	t.Run("ListAllError", func(t *testing.T) {
		_, err := svc.ListAll()
		if err == nil {
			t.Error("Expected error from ListAll")
		}
	})

	t.Run("CreateColorError", func(t *testing.T) {
		_, err := svc.CreateColor("Fail", true, "user")
		if err == nil {
			t.Error("Expected error from CreateColor")
		}
	})

	t.Run("UpdateColorError", func(t *testing.T) {
		_, err := svc.UpdateColor(1, "Fail", true, "user")
		if err == nil {
			t.Error("Expected error from UpdateColor")
		}
	})

	t.Run("DeleteColorError", func(t *testing.T) {
		err := svc.DeleteColor(1)
		if err == nil {
			t.Error("Expected error from DeleteColor")
		}
	})

	t.Run("UpdateColorExecError", func(t *testing.T) {
		// Mock where GetByID succeeds but Update fails
		failRepo := &MockVehicleColorRepository{
			GetByIDFn: func(id int64) (*VehicleColor, error) {
				return &VehicleColor{ID: 1}, nil
			},
			UpdateFn: func(color *VehicleColor) error {
				return errors.New("update failed")
			},
		}
		failSvc := NewVehicleColorService(failRepo)
		_, err := failSvc.UpdateColor(1, "Name", true, "user")
		if err == nil {
			t.Error("Expected error from UpdateColor when Update fails")
		}
	})
}

func TestListPaged(t *testing.T) {
	mockRepo := &MockVehicleColorRepository{
		GetPagedFn: func(limit, offset int, sortBy, sortOrder string) ([]VehicleColor, error) {
			return []VehicleColor{{ID: 1, Name: "Red"}}, nil
		},
		CountFn: func() (int, error) {
			return 1, nil
		},
	}

	svc := NewVehicleColorService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		colors, total, err := svc.ListPaged(1, 10, "id", "desc")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if total != 1 || len(colors) != 1 {
			t.Errorf("Expected total 1 and 1 color, got %d and %d", total, len(colors))
		}
	})

	t.Run("Defaults", func(t *testing.T) {
		// Test negative page/pageSize
		colors, _, err := svc.ListPaged(-1, -1, "id", "desc")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(colors) != 1 {
			t.Errorf("Expected 1 color, got %d", len(colors))
		}
	})

	t.Run("RepoError", func(t *testing.T) {
		failRepo := &MockVehicleColorRepository{
			GetPagedFn: func(limit, offset int, sortBy, sortOrder string) ([]VehicleColor, error) {
				return nil, errors.New("fail")
			},
		}
		failSvc := NewVehicleColorService(failRepo)
		_, _, err := failSvc.ListPaged(1, 10, "id", "desc")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("CountError", func(t *testing.T) {
		failRepo := &MockVehicleColorRepository{
			GetPagedFn: func(limit, offset int, sortBy, sortOrder string) ([]VehicleColor, error) {
				return []VehicleColor{}, nil
			},
			CountFn: func() (int, error) {
				return 0, errors.New("fail")
			},
		}
		failSvc := NewVehicleColorService(failRepo)
		_, _, err := failSvc.ListPaged(1, 10, "id", "desc")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
