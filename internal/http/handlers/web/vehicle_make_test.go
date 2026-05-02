package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang-htmx-bulma/internal/model"
	"golang-htmx-bulma/internal/view"

	"github.com/gin-gonic/gin"
)

func setupMakeTestRouter(svc *MockVehicleMakeService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.HTMLRender = view.NewRenderer("../../../../templates")
	return r
}

func TestListMake(t *testing.T) {
	mockSvc := &MockVehicleMakeService{
		ListAllFn: func() ([]model.VehicleMake, error) {
			return []model.VehicleMake{{ID: 1, Name: "Toyota"}}, nil
		},
	}

	h := NewVehicleMakeHandler(mockSvc)
	r := setupMakeTestRouter(mockSvc)
	r.GET("/vehicle-makes", h.List)

	req, _ := http.NewRequest("GET", "/vehicle-makes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCreateFormMake(t *testing.T) {
	mockSvc := &MockVehicleMakeService{}
	h := NewVehicleMakeHandler(mockSvc)
	r := setupMakeTestRouter(mockSvc)
	r.GET("/vehicle-makes/new", h.CreateForm)

	req, _ := http.NewRequest("GET", "/vehicle-makes/new", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCreateMakeHandler(t *testing.T) {
	mockSvc := &MockVehicleMakeService{
		CreateMakeFn: func(name string, status bool, user string) (*model.VehicleMake, error) {
			return &model.VehicleMake{ID: 1, Name: name}, nil
		},
	}
	h := NewVehicleMakeHandler(mockSvc)
	r := setupMakeTestRouter(mockSvc)
	r.POST("/vehicle-makes", h.Create)

	req, _ := http.NewRequest("POST", "/vehicle-makes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected redirect, got %d", w.Code)
	}
}

func TestEditFormMake(t *testing.T) {
	mockSvc := &MockVehicleMakeService{
		FindByIDFn: func(id int64) (*model.VehicleMake, error) {
			return &model.VehicleMake{ID: id, Name: "Toyota"}, nil
		},
	}
	h := NewVehicleMakeHandler(mockSvc)
	r := setupMakeTestRouter(mockSvc)
	r.GET("/vehicle-makes/:id/edit", h.EditForm)

	req, _ := http.NewRequest("GET", "/vehicle-makes/1/edit", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestUpdateMakeHandler(t *testing.T) {
	mockSvc := &MockVehicleMakeService{
		UpdateMakeFn: func(id int64, name string, status bool, user string) (*model.VehicleMake, error) {
			return &model.VehicleMake{ID: id, Name: name}, nil
		},
	}
	h := NewVehicleMakeHandler(mockSvc)
	r := setupMakeTestRouter(mockSvc)
	r.POST("/vehicle-makes/:id", h.Update)

	req, _ := http.NewRequest("POST", "/vehicle-makes/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected redirect, got %d", w.Code)
	}
}

func TestDeleteMakeHandler(t *testing.T) {
	mockSvc := &MockVehicleMakeService{
		DeleteMakeFn: func(id int64) error {
			return nil
		},
	}
	h := NewVehicleMakeHandler(mockSvc)
	r := setupMakeTestRouter(mockSvc)
	r.DELETE("/vehicle-makes/:id", h.Delete)

	req, _ := http.NewRequest("DELETE", "/vehicle-makes/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandlerErrorsMake(t *testing.T) {
	errSvc := &MockVehicleMakeService{
		ListAllFn: func() ([]model.VehicleMake, error) {
			return nil, errors.New("fail")
		},
		FindByIDFn: func(id int64) (*model.VehicleMake, error) {
			return nil, errors.New("fail")
		},
		CreateMakeFn: func(name string, status bool, user string) (*model.VehicleMake, error) {
			return nil, errors.New("fail")
		},
		UpdateMakeFn: func(id int64, name string, status bool, user string) (*model.VehicleMake, error) {
			return nil, errors.New("fail")
		},
		DeleteMakeFn: func(id int64) error {
			return errors.New("fail")
		},
	}

	h := NewVehicleMakeHandler(errSvc)
	r := setupMakeTestRouter(errSvc)

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
}

func TestDeleteConfirmMake(t *testing.T) {
	mockSvc := &MockVehicleMakeService{
		FindByIDFn: func(id int64) (*model.VehicleMake, error) {
			return &model.VehicleMake{ID: id, Name: "Toyota"}, nil
		},
	}
	h := NewVehicleMakeHandler(mockSvc)
	r := setupMakeTestRouter(mockSvc)
	r.GET("/vehicle-makes/:id/delete", h.DeleteConfirm)

	req, _ := http.NewRequest("GET", "/vehicle-makes/1/delete", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestNotFoundErrorsMake(t *testing.T) {
	// A service that returns nil (not found)
	nilSvc := &MockVehicleMakeService{
		FindByIDFn: func(id int64) (*model.VehicleMake, error) {
			return nil, nil
		},
	}
	h := NewVehicleMakeHandler(nilSvc)
	r := setupMakeTestRouter(nilSvc)
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
}
