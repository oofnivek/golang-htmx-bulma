package service

import (
	"errors"
	"golang-htmx-bulma/internal/model"
	"testing"
)

func TestListAllRole(t *testing.T) {
	mockRepo := &MockRoleRepository{
		GetAllFn: func() ([]model.Role, error) {
			return []model.Role{{ID: "ADMIN"}}, nil
		},
	}
	svc := NewRoleService(mockRepo)
	roles, err := svc.ListAll()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(roles) != 1 {
		t.Errorf("Expected 1 role, got %d", len(roles))
	}
}

func TestListPagedRole(t *testing.T) {
	mockRepo := &MockRoleRepository{
		GetPagedFn: func(limit, offset int, sortBy, sortOrder, search string) ([]model.Role, error) {
			return []model.Role{{ID: "ADMIN", Name: "Administrator"}}, nil
		},
		CountFn: func(search string) (int, error) {
			return 1, nil
		},
	}

	svc := NewRoleService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		roles, total, err := svc.ListPaged(1, 10, "id", "asc", "")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if total != 1 || len(roles) != 1 {
			t.Errorf("Expected total 1 and 1 role, got %d and %d", total, len(roles))
		}
	})

	t.Run("Defaults", func(t *testing.T) {
		roles, _, err := svc.ListPaged(0, 0, "id", "asc", "")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(roles) != 1 {
			t.Errorf("Expected 1 role, got %d", len(roles))
		}
	})

	t.Run("GetPagedError", func(t *testing.T) {
		failRepo := &MockRoleRepository{
			GetPagedFn: func(limit, offset int, sortBy, sortOrder, search string) ([]model.Role, error) {
				return nil, errors.New("fail")
			},
		}
		failSvc := NewRoleService(failRepo)
		_, _, err := failSvc.ListPaged(1, 10, "id", "asc", "")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("CountError", func(t *testing.T) {
		failRepo := &MockRoleRepository{
			GetPagedFn: func(limit, offset int, sortBy, sortOrder, search string) ([]model.Role, error) {
				return []model.Role{}, nil
			},
			CountFn: func(search string) (int, error) {
				return 0, errors.New("fail")
			},
		}
		failSvc := NewRoleService(failRepo)
		_, _, err := failSvc.ListPaged(1, 10, "id", "asc", "")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestFindByIDRole(t *testing.T) {
	mockRepo := &MockRoleRepository{
		GetByIDFn: func(id string) (*model.Role, error) {
			if id == "ADMIN" {
				return &model.Role{ID: "ADMIN"}, nil
			}
			return nil, nil
		},
	}
	svc := NewRoleService(mockRepo)

	role, err := svc.FindByID("ADMIN")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if role == nil {
		t.Fatal("Expected role, got nil")
	}
}

func TestCreateRole(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockRoleRepository{
			CreateFn: func(role *model.Role) error {
				return nil
			},
		}
		svc := NewRoleService(mockRepo)
		role, err := svc.CreateRole("ADMIN", "Admin")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if role.ID != "ADMIN" {
			t.Errorf("Expected ADMIN, got %s", role.ID)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo := &MockRoleRepository{
			CreateFn: func(role *model.Role) error {
				return errors.New("fail")
			},
		}
		svc := NewRoleService(mockRepo)
		_, err := svc.CreateRole("ADMIN", "Admin")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestUpdateRole(t *testing.T) {
	mockRepo := &MockRoleRepository{
		GetByIDFn: func(id string) (*model.Role, error) {
			if id == "ADMIN" {
				return &model.Role{ID: "ADMIN", Name: "Old"}, nil
			}
			return nil, nil
		},
		UpdateFn: func(role *model.Role) error {
			return nil
		},
	}
	svc := NewRoleService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		role, err := svc.UpdateRole("ADMIN", "New")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if role.Name != "New" {
			t.Errorf("Expected New, got %s", role.Name)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		role, err := svc.UpdateRole("NONE", "New")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if role != nil {
			t.Error("Expected nil, got role")
		}
	})

	t.Run("GetError", func(t *testing.T) {
		errRepo := &MockRoleRepository{
			GetByIDFn: func(id string) (*model.Role, error) {
				return nil, errors.New("fail")
			},
		}
		errSvc := NewRoleService(errRepo)
		_, err := errSvc.UpdateRole("ADMIN", "New")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("UpdateError", func(t *testing.T) {
		errRepo := &MockRoleRepository{
			GetByIDFn: func(id string) (*model.Role, error) {
				return &model.Role{ID: "ADMIN"}, nil
			},
			UpdateFn: func(role *model.Role) error {
				return errors.New("fail")
			},
		}
		errSvc := NewRoleService(errRepo)
		_, err := errSvc.UpdateRole("ADMIN", "New")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestDeleteRole(t *testing.T) {
	mockRepo := &MockRoleRepository{
		DeleteFn: func(id string) error {
			return nil
		},
	}
	svc := NewRoleService(mockRepo)
	err := svc.DeleteRole("ADMIN")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
