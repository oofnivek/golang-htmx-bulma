package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang-htmx-bulma/internal/vehicle"
	"golang-htmx-bulma/internal/view"

	"github.com/gin-gonic/gin"
)

func setupVehicleTypeTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.HTMLRender = view.NewRenderer("../../../../templates")
	return r
}

func TestRegisterVehicleTypeRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := &MockVehicleTypeService{}
	group := r.Group("/")
	RegisterVehicleTypeRoutes(group, svc)
}

func TestVehicleTypeList(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := &MockVehicleTypeService{
			ListPagedFn: func(page, pageSize int, sortBy, sortOrder string) ([]vehicle.VehicleType, int, error) {
				return []vehicle.VehicleType{{ID: 1, Name: "Sedan"}}, 1, nil
			},
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.GET("/vehicle-types", h.List)

		req, _ := http.NewRequest("GET", "/vehicle-types", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		svc := &MockVehicleTypeService{
			ListPagedFn: func(page, pageSize int, sortBy, sortOrder string) ([]vehicle.VehicleType, int, error) {
				return nil, 0, errors.New("fail")
			},
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.GET("/vehicle-types", h.List)

		req, _ := http.NewRequest("GET", "/vehicle-types", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestVehicleTypeCreateForm(t *testing.T) {
	svc := &MockVehicleTypeService{}
	h := NewVehicleTypeHandler(svc)
	r := setupVehicleTypeTestRouter()
	r.GET("/vehicle-types/new", h.CreateForm)

	req, _ := http.NewRequest("GET", "/vehicle-types/new", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestVehicleTypeCreate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var capturedUser string
		svc := &MockVehicleTypeService{
			CreateFn: func(name string, status bool, user string) (*vehicle.VehicleType, error) {
				capturedUser = user
				return &vehicle.VehicleType{ID: 1, Name: name}, nil
			},
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.Use(func(c *gin.Context) {
			c.Set("user_email", "tester@example.com")
			c.Next()
		})
		r.POST("/vehicle-types", h.Create)

		req, _ := http.NewRequest("POST", "/vehicle-types", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected redirect, got %d", w.Code)
		}
		if capturedUser != "tester@example.com" {
			t.Errorf("Expected user email, got %s", capturedUser)
		}
	})

	t.Run("Error", func(t *testing.T) {
		svc := &MockVehicleTypeService{
			CreateFn: func(name string, status bool, user string) (*vehicle.VehicleType, error) {
				return nil, errors.New("fail")
			},
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.POST("/vehicle-types", h.Create)

		req, _ := http.NewRequest("POST", "/vehicle-types", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestVehicleTypeEditForm(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		svc := &MockVehicleTypeService{
			FindByIDFn: func(id int64) (*vehicle.VehicleType, error) {
				return &vehicle.VehicleType{ID: id, Name: "Sedan"}, nil
			},
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.GET("/vehicle-types/:id/edit", h.EditForm)

		req, _ := http.NewRequest("GET", "/vehicle-types/1/edit", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		svc := &MockVehicleTypeService{
			FindByIDFn: func(id int64) (*vehicle.VehicleType, error) { return nil, nil },
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.GET("/vehicle-types/:id/edit", h.EditForm)

		req, _ := http.NewRequest("GET", "/vehicle-types/99/edit", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		svc := &MockVehicleTypeService{
			FindByIDFn: func(id int64) (*vehicle.VehicleType, error) { return nil, errors.New("fail") },
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.GET("/vehicle-types/:id/edit", h.EditForm)

		req, _ := http.NewRequest("GET", "/vehicle-types/1/edit", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestVehicleTypeUpdate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var capturedUser string
		svc := &MockVehicleTypeService{
			UpdateFn: func(id int64, name string, status bool, user string) (*vehicle.VehicleType, error) {
				capturedUser = user
				return &vehicle.VehicleType{ID: id, Name: name}, nil
			},
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.Use(func(c *gin.Context) {
			c.Set("user_email", "tester@example.com")
			c.Next()
		})
		r.POST("/vehicle-types/:id", h.Update)

		req, _ := http.NewRequest("POST", "/vehicle-types/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected redirect, got %d", w.Code)
		}
		if capturedUser != "tester@example.com" {
			t.Errorf("Expected user email, got %s", capturedUser)
		}
	})

	t.Run("Error", func(t *testing.T) {
		svc := &MockVehicleTypeService{
			UpdateFn: func(id int64, name string, status bool, user string) (*vehicle.VehicleType, error) {
				return nil, errors.New("fail")
			},
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.POST("/vehicle-types/:id", h.Update)

		req, _ := http.NewRequest("POST", "/vehicle-types/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestVehicleTypeDelete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := &MockVehicleTypeService{
			DeleteFn: func(id int64) error { return nil },
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.POST("/vehicle-types/:id/delete", h.Delete)

		req, _ := http.NewRequest("POST", "/vehicle-types/1/delete", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		svc := &MockVehicleTypeService{
			DeleteFn: func(id int64) error { return errors.New("fail") },
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.POST("/vehicle-types/:id/delete", h.Delete)

		req, _ := http.NewRequest("POST", "/vehicle-types/1/delete", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestVehicleTypeDeleteConfirm(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		svc := &MockVehicleTypeService{
			FindByIDFn: func(id int64) (*vehicle.VehicleType, error) {
				return &vehicle.VehicleType{ID: id, Name: "Sedan"}, nil
			},
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.GET("/vehicle-types/:id/delete", h.DeleteConfirm)

		req, _ := http.NewRequest("GET", "/vehicle-types/1/delete", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		svc := &MockVehicleTypeService{
			FindByIDFn: func(id int64) (*vehicle.VehicleType, error) { return nil, nil },
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.GET("/vehicle-types/:id/delete", h.DeleteConfirm)

		req, _ := http.NewRequest("GET", "/vehicle-types/99/delete", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		svc := &MockVehicleTypeService{
			FindByIDFn: func(id int64) (*vehicle.VehicleType, error) { return nil, errors.New("fail") },
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.GET("/vehicle-types/:id/delete", h.DeleteConfirm)

		req, _ := http.NewRequest("GET", "/vehicle-types/1/delete", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}

func TestVehicleTypeView(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		svc := &MockVehicleTypeService{
			FindByIDFn: func(id int64) (*vehicle.VehicleType, error) {
				return &vehicle.VehicleType{ID: id, Name: "Sedan"}, nil
			},
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.GET("/vehicle-types/:id", h.View)

		req, _ := http.NewRequest("GET", "/vehicle-types/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		svc := &MockVehicleTypeService{
			FindByIDFn: func(id int64) (*vehicle.VehicleType, error) { return nil, nil },
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.GET("/vehicle-types/:id", h.View)

		req, _ := http.NewRequest("GET", "/vehicle-types/99", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", w.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		svc := &MockVehicleTypeService{
			FindByIDFn: func(id int64) (*vehicle.VehicleType, error) { return nil, errors.New("fail") },
		}
		h := NewVehicleTypeHandler(svc)
		r := setupVehicleTypeTestRouter()
		r.GET("/vehicle-types/:id", h.View)

		req, _ := http.NewRequest("GET", "/vehicle-types/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500, got %d", w.Code)
		}
	})
}
