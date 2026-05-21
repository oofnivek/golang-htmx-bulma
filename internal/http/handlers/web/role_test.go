package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang-htmx-bulma/internal/user"
	"golang-htmx-bulma/internal/view"

	"github.com/gin-gonic/gin"
)

func setupRoleTestRouter(svc *MockRoleService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.HTMLRender = view.NewRenderer("../../../../templates")
	return r
}

func TestListRole(t *testing.T) {
	mockSvc := &MockRoleService{
		ListPagedFn: func(page, pageSize int, sortBy, sortOrder, search string) ([]user.Role, int, error) {
			return []user.Role{{ID: "ADMIN", Name: "Administrator"}}, 1, nil
		},
	}

	h := NewRoleHandler(mockSvc)
	r := setupRoleTestRouter(mockSvc)
	r.GET("/roles", h.List)

	t.Run("Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/roles", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		errSvc := &MockRoleService{
			ListPagedFn: func(page, pageSize int, sortBy, sortOrder, search string) ([]user.Role, int, error) {
				return nil, 0, errors.New("fail")
			},
		}
		h2 := NewRoleHandler(errSvc)
		r2 := setupRoleTestRouter(errSvc)
		r2.GET("/roles", h2.List)
		req, _ := http.NewRequest("GET", "/roles", nil)
		w := httptest.NewRecorder()
		r2.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestCreateFormRole(t *testing.T) {
	mockSvc := &MockRoleService{}
	h := NewRoleHandler(mockSvc)
	r := setupRoleTestRouter(mockSvc)
	r.GET("/roles/new", h.CreateForm)

	req, _ := http.NewRequest("GET", "/roles/new", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestCreateRoleHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := &MockRoleService{
			CreateRoleFn: func(id, name string) (*user.Role, error) {
				return &user.Role{ID: id, Name: name}, nil
			},
		}
		h := NewRoleHandler(mockSvc)
		r := setupRoleTestRouter(mockSvc)
		r.POST("/roles", h.Create)

		req, _ := http.NewRequest("POST", "/roles", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected redirect, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockSvc := &MockRoleService{
			CreateRoleFn: func(id, name string) (*user.Role, error) {
				return nil, errors.New("fail")
			},
		}
		h := NewRoleHandler(mockSvc)
		r := setupRoleTestRouter(mockSvc)
		r.POST("/roles", h.Create)

		req, _ := http.NewRequest("POST", "/roles", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestEditFormRole(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		mockSvc := &MockRoleService{
			FindByIDFn: func(id string) (*user.Role, error) {
				return &user.Role{ID: id, Name: "Admin"}, nil
			},
		}
		h := NewRoleHandler(mockSvc)
		r := setupRoleTestRouter(mockSvc)
		r.GET("/roles/:id/edit", h.EditForm)

		req, _ := http.NewRequest("GET", "/roles/ADMIN/edit", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockSvc := &MockRoleService{
			FindByIDFn: func(id string) (*user.Role, error) {
				return nil, nil
			},
		}
		h := NewRoleHandler(mockSvc)
		r := setupRoleTestRouter(mockSvc)
		r.GET("/roles/:id/edit", h.EditForm)

		req, _ := http.NewRequest("GET", "/roles/NONE/edit", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockSvc := &MockRoleService{
			FindByIDFn: func(id string) (*user.Role, error) {
				return nil, errors.New("fail")
			},
		}
		h := NewRoleHandler(mockSvc)
		r := setupRoleTestRouter(mockSvc)
		r.GET("/roles/:id/edit", h.EditForm)

		req, _ := http.NewRequest("GET", "/roles/FAIL/edit", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestUpdateRoleHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := &MockRoleService{
			UpdateRoleFn: func(id, name string) (*user.Role, error) {
				return &user.Role{ID: id, Name: name}, nil
			},
		}
		h := NewRoleHandler(mockSvc)
		r := setupRoleTestRouter(mockSvc)
		r.POST("/roles/:id", h.Update)

		req, _ := http.NewRequest("POST", "/roles/ADMIN", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected redirect, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockSvc := &MockRoleService{
			UpdateRoleFn: func(id, name string) (*user.Role, error) {
				return nil, errors.New("fail")
			},
		}
		h := NewRoleHandler(mockSvc)
		r := setupRoleTestRouter(mockSvc)
		r.POST("/roles/:id", h.Update)

		req, _ := http.NewRequest("POST", "/roles/ADMIN", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestDeleteRoleHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := &MockRoleService{
			DeleteRoleFn: func(id string) error {
				return nil
			},
		}
		h := NewRoleHandler(mockSvc)
		r := setupRoleTestRouter(mockSvc)
		r.DELETE("/roles/:id", h.Delete)

		req, _ := http.NewRequest("DELETE", "/roles/ADMIN", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockSvc := &MockRoleService{
			DeleteRoleFn: func(id string) error {
				return errors.New("fail")
			},
		}
		h := NewRoleHandler(mockSvc)
		r := setupRoleTestRouter(mockSvc)
		r.DELETE("/roles/:id", h.Delete)

		req, _ := http.NewRequest("DELETE", "/roles/ADMIN", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestDeleteConfirmRole(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		mockSvc := &MockRoleService{
			FindByIDFn: func(id string) (*user.Role, error) {
				return &user.Role{ID: id, Name: "Admin"}, nil
			},
		}
		h := NewRoleHandler(mockSvc)
		r := setupRoleTestRouter(mockSvc)
		r.GET("/roles/:id/delete", h.DeleteConfirm)

		req, _ := http.NewRequest("GET", "/roles/ADMIN/delete", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockSvc := &MockRoleService{
			FindByIDFn: func(id string) (*user.Role, error) {
				return nil, nil
			},
		}
		h := NewRoleHandler(mockSvc)
		r := setupRoleTestRouter(mockSvc)
		r.GET("/roles/:id/delete", h.DeleteConfirm)

		req, _ := http.NewRequest("GET", "/roles/NONE/delete", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockSvc := &MockRoleService{
			FindByIDFn: func(id string) (*user.Role, error) {
				return nil, errors.New("fail")
			},
		}
		h := NewRoleHandler(mockSvc)
		r := setupRoleTestRouter(mockSvc)
		r.GET("/roles/:id/delete", h.DeleteConfirm)

		req, _ := http.NewRequest("GET", "/roles/FAIL/delete", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}
