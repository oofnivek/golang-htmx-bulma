package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang-htmx-bulma/internal/model"
	"golang-htmx-bulma/internal/service"
	"golang-htmx-bulma/internal/view"

	"github.com/gin-gonic/gin"
)

func setupTestRouter(svc service.VehicleColorService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	// Initialize renderer pointing to the templates directory
	// In tests, we need to go up to the project root
	r.HTMLRender = view.NewRenderer("../../../../templates")
	return r
}

func TestList(t *testing.T) {
	mockSvc := &MockVehicleColorService{
		ListPagedFn: func(page, pageSize int, sortBy, sortOrder, search string) ([]model.VehicleColor, int, error) {
			return []model.VehicleColor{{ID: 1, Name: "Red"}}, 1, nil
		},
	}

	h := NewVehicleColorHandler(mockSvc)
	r := setupTestRouter(mockSvc)
	r.GET("/vehicle-colors", h.List)

	req, _ := http.NewRequest("GET", "/vehicle-colors", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCreateForm(t *testing.T) {
	mockSvc := &MockVehicleColorService{}
	h := NewVehicleColorHandler(mockSvc)
	r := setupTestRouter(mockSvc)
	r.GET("/vehicle-colors/new", h.CreateForm)

	req, _ := http.NewRequest("GET", "/vehicle-colors/new", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteConfirm(t *testing.T) {
	mockSvc := &MockVehicleColorService{
		FindByIDFn: func(id int64) (*model.VehicleColor, error) {
			return &model.VehicleColor{ID: id, Name: "Red"}, nil
		},
	}
	h := NewVehicleColorHandler(mockSvc)
	r := setupTestRouter(mockSvc)
	r.GET("/vehicle-colors/:id/delete", h.DeleteConfirm)

	req, _ := http.NewRequest("GET", "/vehicle-colors/1/delete", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCreate(t *testing.T) {
	mockSvc := &MockVehicleColorService{
		CreateColorFn: func(name string, status bool, user string) (*model.VehicleColor, error) {
			return &model.VehicleColor{ID: 1, Name: name}, nil
		},
	}
	h := NewVehicleColorHandler(mockSvc)
	r := setupTestRouter(mockSvc)
	r.POST("/vehicle-colors", h.Create)

	req, _ := http.NewRequest("POST", "/vehicle-colors", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected redirect, got %d", w.Code)
	}
}

func TestEditForm(t *testing.T) {
	mockSvc := &MockVehicleColorService{
		FindByIDFn: func(id int64) (*model.VehicleColor, error) {
			return &model.VehicleColor{ID: id, Name: "Red"}, nil
		},
	}
	h := NewVehicleColorHandler(mockSvc)
	r := setupTestRouter(mockSvc)
	r.GET("/vehicle-colors/:id/edit", h.EditForm)

	req, _ := http.NewRequest("GET", "/vehicle-colors/1/edit", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestUpdate(t *testing.T) {
	mockSvc := &MockVehicleColorService{
		UpdateColorFn: func(id int64, name string, status bool, user string) (*model.VehicleColor, error) {
			return &model.VehicleColor{ID: id, Name: name}, nil
		},
	}
	h := NewVehicleColorHandler(mockSvc)
	r := setupTestRouter(mockSvc)
	r.POST("/vehicle-colors/:id", h.Update)

	req, _ := http.NewRequest("POST", "/vehicle-colors/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected redirect, got %d", w.Code)
	}
}

func TestDelete(t *testing.T) {
	mockSvc := &MockVehicleColorService{
		DeleteColorFn: func(id int64) error {
			return nil
		},
	}
	h := NewVehicleColorHandler(mockSvc)
	r := setupTestRouter(mockSvc)
	r.DELETE("/vehicle-colors/:id", h.Delete)

	req, _ := http.NewRequest("DELETE", "/vehicle-colors/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandlerErrors(t *testing.T) {
	// A service that always returns errors
	errSvc := &MockVehicleColorService{
		ListPagedFn: func(page, pageSize int, sortBy, sortOrder, search string) ([]model.VehicleColor, int, error) {
			return nil, 0, errors.New("fail")
		},
		FindByIDFn: func(id int64) (*model.VehicleColor, error) {
			return nil, errors.New("fail")
		},
		CreateColorFn: func(name string, status bool, user string) (*model.VehicleColor, error) {
			return nil, errors.New("fail")
		},
		UpdateColorFn: func(id int64, name string, status bool, user string) (*model.VehicleColor, error) {
			return nil, errors.New("fail")
		},
		DeleteColorFn: func(id int64) error {
			return errors.New("fail")
		},
	}

	h := NewVehicleColorHandler(errSvc)
	r := setupTestRouter(errSvc)

	r.GET("/list", h.List)
	r.POST("/create", h.Create)
	r.GET("/edit/:id", h.EditForm)
	r.POST("/update/:id", h.Update)
	r.DELETE("/delete/:id", h.Delete)
	r.GET("/delete-confirm/:id", h.DeleteConfirm)

	t.Run("ListError", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/list", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})

	t.Run("CreateError", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/create", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})

	t.Run("EditFormError", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/edit/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})

	t.Run("UpdateError", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/update/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})

	t.Run("DeleteError", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/delete/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})

	t.Run("DeleteConfirmError", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/delete-confirm/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})

	t.Run("NotFoundErrors", func(t *testing.T) {
		// A service that returns nil (not found)
		nilSvc := &MockVehicleColorService{
			FindByIDFn: func(id int64) (*model.VehicleColor, error) {
				return nil, nil
			},
		}
		h := NewVehicleColorHandler(nilSvc)
		r := setupTestRouter(nilSvc)
		r.GET("/edit/:id", h.EditForm)
		r.GET("/delete-confirm/:id", h.DeleteConfirm)

		t.Run("EditFormNotFound", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/edit/99", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != http.StatusNotFound {
				t.Errorf("Expected 404, got %d", w.Code)
			}
		})

		t.Run("DeleteConfirmNotFound", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/delete-confirm/99", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != http.StatusNotFound {
				t.Errorf("Expected 404, got %d", w.Code)
			}
		})
	})
}
