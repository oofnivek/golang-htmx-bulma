package service

import (
	"errors"
	"golang-htmx-bulma/internal/model"
	"testing"
)

func TestCreateMake(t *testing.T) {
	mockRepo := &MockVehicleMakeRepository{
		CreateFn: func(make *model.VehicleMake) error {
			make.ID = 1
			return nil
		},
	}

	svc := NewVehicleMakeService(mockRepo)

	make, err := svc.CreateMake("Toyota", true, "test-user")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if make.ID != 1 {
		t.Errorf("Expected ID 1, got %d", make.ID)
	}

	if make.Name != "Toyota" {
		t.Errorf("Expected name 'Toyota', got %s", make.Name)
	}

	if make.UpdatedBy == nil || *make.UpdatedBy != "test-user" {
		t.Errorf("Expected UpdatedBy 'test-user', got %v", make.UpdatedBy)
	}
}

func TestFindByIDMake(t *testing.T) {
	mockRepo := &MockVehicleMakeRepository{
		GetByIDFn: func(id int64) (*model.VehicleMake, error) {
			if id == 1 {
				return &model.VehicleMake{ID: 1, Name: "Honda"}, nil
			}
			return nil, nil
		},
	}

	svc := NewVehicleMakeService(mockRepo)

	t.Run("Found", func(t *testing.T) {
		make, err := svc.FindByID(1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if make == nil || make.Name != "Honda" {
			t.Errorf("Expected 'Honda', got %v", make)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		make, err := svc.FindByID(99)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if make != nil {
			t.Errorf("Expected nil, got %v", make)
		}
	})
}

func TestListAllMake(t *testing.T) {
	mockRepo := &MockVehicleMakeRepository{
		GetAllFn: func() ([]model.VehicleMake, error) {
			return []model.VehicleMake{
				{ID: 1, Name: "Toyota"},
				{ID: 2, Name: "Honda"},
			}, nil
		},
	}

	svc := NewVehicleMakeService(mockRepo)
	makes, err := svc.ListAll()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(makes) != 2 {
		t.Errorf("Expected 2 makes, got %d", len(makes))
	}
}

func TestUpdateMake(t *testing.T) {
	mockRepo := &MockVehicleMakeRepository{
		GetByIDFn: func(id int64) (*model.VehicleMake, error) {
			if id == 1 {
				return &model.VehicleMake{ID: 1, Name: "Old Name"}, nil
			}
			return nil, nil
		},
		UpdateFn: func(make *model.VehicleMake) error {
			return nil
		},
	}

	svc := NewVehicleMakeService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		make, err := svc.UpdateMake(1, "New Name", false, "admin")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if make.Name != "New Name" {
			t.Errorf("Expected 'New Name', got %s", make.Name)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		make, err := svc.UpdateMake(99, "New Name", false, "admin")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if make != nil {
			t.Errorf("Expected nil for non-existent make")
		}
	})
}

func TestDeleteMake(t *testing.T) {
	mockRepo := &MockVehicleMakeRepository{
		DeleteFn: func(id int64) error {
			return nil
		},
	}

	svc := NewVehicleMakeService(mockRepo)
	err := svc.DeleteMake(1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestServiceErrorsMake(t *testing.T) {
	errRepo := &MockVehicleMakeRepository{
		GetAllFn: func() ([]model.VehicleMake, error) {
			return nil, errors.New("db error")
		},
		CreateFn: func(make *model.VehicleMake) error {
			return errors.New("db error")
		},
		UpdateFn: func(make *model.VehicleMake) error {
			return errors.New("db error")
		},
		GetByIDFn: func(id int64) (*model.VehicleMake, error) {
			return nil, errors.New("db error")
		},
		DeleteFn: func(id int64) error {
			return errors.New("db error")
		},
	}

	svc := NewVehicleMakeService(errRepo)

	t.Run("ListAllError", func(t *testing.T) {
		_, err := svc.ListAll()
		if err == nil {
			t.Error("Expected error from ListAll")
		}
	})

	t.Run("CreateMakeError", func(t *testing.T) {
		_, err := svc.CreateMake("Fail", true, "user")
		if err == nil {
			t.Error("Expected error from CreateMake")
		}
	})

	t.Run("UpdateMakeError", func(t *testing.T) {
		_, err := svc.UpdateMake(1, "Fail", true, "user")
		if err == nil {
			t.Error("Expected error from UpdateMake")
		}
	})

	t.Run("DeleteMakeError", func(t *testing.T) {
		err := svc.DeleteMake(1)
		if err == nil {
			t.Error("Expected error from DeleteMake")
		}
	})

	t.Run("UpdateMakeExecError", func(t *testing.T) {
		// Mock where GetByID succeeds but Update fails
		failRepo := &MockVehicleMakeRepository{
			GetByIDFn: func(id int64) (*model.VehicleMake, error) {
				return &model.VehicleMake{ID: 1}, nil
			},
			UpdateFn: func(make *model.VehicleMake) error {
				return errors.New("update failed")
			},
		}
		failSvc := NewVehicleMakeService(failRepo)
		_, err := failSvc.UpdateMake(1, "Name", true, "user")
		if err == nil {
			t.Error("Expected error from UpdateMake when Update fails")
		}
	})
}
